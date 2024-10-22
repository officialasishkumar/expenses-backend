package handlers

import (
	"context"
	"encoding/csv"
	"expenses-backend/db"
	"expenses-backend/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson" // Added bson import
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BalanceSheetRow represents a row in the balance sheet
type BalanceSheetRow struct {
    Name         string  `json:"name"`
    Email        string  `json:"email"`
    MobileNumber string  `json:"mobile_number"`
    TotalSpent   float64 `json:"total_spent"`
    TotalOwed    float64 `json:"total_owed"`
    NetBalance   float64 `json:"net_balance"`
}

// DownloadBalanceSheet generates and sends a CSV balance sheet
func DownloadBalanceSheet(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Fetch all users
    cursor, err := db.UsersCol.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    balanceRows := []BalanceSheetRow{}

    for cursor.Next(ctx) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
            return
        }

        // Calculate total spent
        totalSpent, err := calculateTotalSpent(ctx, user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total spent"})
            return
        }

        // Calculate total owed
        totalOwed, err := calculateTotalOwed(ctx, user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total owed"})
            return
        }

        netBalance := totalSpent - totalOwed

        balanceRows = append(balanceRows, BalanceSheetRow{
            Name:         user.Name,
            Email:        user.Email,
            MobileNumber: user.MobileNumber,
            TotalSpent:   totalSpent,
            TotalOwed:    totalOwed,
            NetBalance:   netBalance,
        })
    }

    // Prepare CSV data
    csvData := [][]string{
        {"Name", "Email", "Mobile Number", "Total Spent", "Total Owed", "Net Balance"},
    }

    for _, r := range balanceRows {
        csvData = append(csvData, []string{
            r.Name,
            r.Email,
            r.MobileNumber,
            fmt.Sprintf("%.2f", r.TotalSpent),
            fmt.Sprintf("%.2f", r.TotalOwed),
            fmt.Sprintf("%.2f", r.NetBalance),
        })
    }

    // Create CSV file in memory
    csvString := &strings.Builder{}
    writer := csv.NewWriter(csvString)
    writer.WriteAll(csvData)
    writer.Flush()

    if err := writer.Error(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSV"})
        return
    }

    // Send CSV as downloadable file
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename=balance_sheet.csv")
    c.Data(http.StatusOK, "text/csv", []byte(csvString.String()))
}

func calculateTotalSpent(ctx context.Context, userID primitive.ObjectID) (float64, error) {
    filter := bson.M{"created_by": userID}
    cursor, err := db.ExpensesCol.Find(ctx, filter)
    if err != nil {
        return 0, err
    }
    defer cursor.Close(ctx)

    totalSpent := 0.0
    for cursor.Next(ctx) {
        var expense models.Expense
        if err := cursor.Decode(&expense); err != nil {
            return 0, err
        }
        totalSpent += expense.Amount
    }
    return totalSpent, nil
}

func calculateTotalOwed(ctx context.Context, userID primitive.ObjectID) (float64, error) {
    filter := bson.M{"participants": userID}
    cursor, err := db.ExpensesCol.Find(ctx, filter)
    if err != nil {
        return 0, err
    }
    defer cursor.Close(ctx)

    totalOwed := 0.0
    for cursor.Next(ctx) {
        var expense models.Expense
        if err := cursor.Decode(&expense); err != nil {
            return 0, err
        }

        // Find the amount owed by the user in split_details
        for key, amount := range expense.SplitDetails {
            if strings.EqualFold(key, userID.Hex()) || strings.EqualFold(key, userID.String()) {
                amt, ok := convertToFloat64(amount)
                if ok {
                    totalOwed += amt
                }
                break
            }
        }
    }
    return totalOwed, nil
}
