package handlers

import (
	"context"
	"expenses-backend/db"
	"expenses-backend/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExpenseInput struct {
	Description  string                 `json:"description" binding:"required"`
	Amount       float64                `json:"amount" binding:"required,gt=0"`
	CreatedBy    string                 `json:"created_by" binding:"required"`
	SplitType    string                 `json:"split_type" binding:"required,oneof=Equal Exact Percentage"`
	Participants []string               `json:"participants" binding:"required,min=1"`
	SplitDetails map[string]interface{} `json:"split_details,omitempty"`
}

var expenseValidate = validator.New()

// AddExpense handles adding a new expense
func AddExpense(c *gin.Context) {
	var input ExpenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if err := expenseValidate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Description = strings.TrimSpace(input.Description)
	input.SplitType = strings.TrimSpace(input.SplitType)

	// Identify creator
	creator, err := identifyUser(input.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'created_by' identifier: " + err.Error()})
		return
	}

	// Identify participants
	participantIDs := []primitive.ObjectID{}
	for _, p := range input.Participants {
		user, err := identifyUser(p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant identifier '" + p + "': " + err.Error()})
			return
		}
		participantIDs = append(participantIDs, user.ID)
	}

	// Validate and compute split
	splits := make(map[primitive.ObjectID]float64)
	switch input.SplitType {
	case "Equal":
		splitAmount := input.Amount / float64(len(participantIDs))
		for _, pid := range participantIDs {
			splits[pid] = splitAmount
		}
	case "Exact":
		if input.SplitDetails == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "split_details required for Exact split"})
			return
		}
		total := 0.0
		for _, v := range input.SplitDetails {
			amount, ok := convertToFloat64(v)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount in split_details"})
				return
			}
			total += amount
		}
		if !almostEqual(total, input.Amount, 0.01) { // Allow small rounding differences
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sum of exact amounts does not equal total amount"})
			return
		}
		for k, v := range input.SplitDetails {
			user, err := identifyUser(k)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant identifier '" + k + "': " + err.Error()})
				return
			}
			amount, ok := convertToFloat64(v)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount for user '" + k + "' in split_details"})
				return
			}
			splits[user.ID] = amount
		}
	case "Percentage":
		if input.SplitDetails == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "split_details required for Percentage split"})
			return
		}
		totalPercent := 0.0
		for _, v := range input.SplitDetails {
			percentage, ok := convertToFloat64(v)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage in split_details"})
				return
			}
			totalPercent += percentage
		}
		if !almostEqual(totalPercent, 100.0, 0.01) { // Allow small rounding differences
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sum of percentages must be exactly 100%"})
			return
		}
		for k, v := range input.SplitDetails {
			user, err := identifyUser(k)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant identifier '" + k + "': " + err.Error()})
				return
			}
			percentage, ok := convertToFloat64(v)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid percentage for user '" + k + "' in split_details"})
				return
			}
			splits[user.ID] = (percentage / 100.0) * input.Amount
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid split_type"})
		return
	}

	// Create Expense
	expense := models.Expense{
		Description:  input.Description,
		Amount:       input.Amount,
		CreatedBy:    creator.ID,
		SplitType:    input.SplitType,
		Participants: participantIDs,
		SplitDetails: make(map[string]interface{}),
		CreatedAt:    time.Now(),
	}

	// Prepare split_details with user identifiers (e.g., email)
	for pid, amt := range splits {
		// Fetch user's email for split_details
		var user models.User
		err := db.UsersCol.FindOne(context.Background(), bson.M{"_id": pid}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch participant details"})
			return
		}
		expense.SplitDetails[user.Email] = amt
	}

	// Insert into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.ExpensesCol.InsertOne(ctx, expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	expense.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, expense)
}

// Helper function to check if two floats are almost equal
func almostEqual(a, b, epsilon float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

// Helper function to convert interface{} to float64
func convertToFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

// GetUserExpenses handles retrieving expenses for a specific user
func GetUserExpenses(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identifier (email, mobile_number, or name) is required"})
		return
	}

	user, err := identifyUser(identifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid identifier: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"created_by": user.ID},
			{"participants": user.ID},
		},
	}

	cursor, err := db.ExpensesCol.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}
	defer cursor.Close(ctx)

	expenses := []models.Expense{}
	for cursor.Next(ctx) {
		var expense models.Expense
		if err := cursor.Decode(&expense); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse expenses"})
			return
		}
		expenses = append(expenses, expense)
	}

	c.JSON(http.StatusOK, expenses)
}

// GetOverallExpenses handles retrieving all expenses
func GetOverallExpenses(c *gin.Context) {
	// Pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := db.ExpensesCol.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		return
	}
	defer cursor.Close(ctx)

	expenses := []models.Expense{}
	for cursor.Next(ctx) {
		var expense models.Expense
		if err := cursor.Decode(&expense); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse expenses"})
			return
		}
		expenses = append(expenses, expense)
	}

	c.JSON(http.StatusOK, expenses)
}
