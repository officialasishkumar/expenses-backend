package handlers

import (
	"context"
	"errors"
	"expenses-backend/db"
	"expenses-backend/models"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

// CreateUser handles creating a new user
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.MobileNumber = strings.TrimSpace(user.MobileNumber)
	user.CreatedAt = time.Now()

	// Insert into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.UsersCol.InsertOne(ctx, user)
	if err != nil {
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, we := range mongoErr.WriteErrors {
				if we.Code == 11000 {
					// Duplicate key error
					c.JSON(http.StatusBadRequest, gin.H{"error": "Email or mobile number already exists"})
					return
				}
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles retrieving user details based on identifier
func GetUser(c *gin.Context) {
	identifier := c.Query("identifier")
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Identifier (email, mobile_number, or name) is required"})
		return
	}

	user, err := identifyUser(identifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// identifyUser identifies a user based on email, phone, or name
func identifyUser(identifier string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	phoneRegex := regexp.MustCompile(`^[6-9]\d{9}$`) // Indian 10-digit phone number

	identifier = strings.TrimSpace(identifier)

	if emailRegex.MatchString(identifier) {
		filter := bson.M{"email": strings.ToLower(identifier)}
		err := db.UsersCol.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return user, errors.New("no user found with the given email")
			}
			return user, err
		}
		return user, nil
	} else if phoneRegex.MatchString(identifier) {
		filter := bson.M{"mobile_number": identifier}
		err := db.UsersCol.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return user, errors.New("no user found with the given mobile number")
			}
			return user, err
		}
		return user, nil
	} else {
		// Treat as name (case-insensitive)
		filter := bson.M{"name": primitive.Regex{Pattern: fmt.Sprintf("^%s$", regexp.QuoteMeta(strings.ToLower(identifier))), Options: "i"}}
		cursor, err := db.UsersCol.Find(ctx, filter)
		if err != nil {
			return user, err
		}
		defer cursor.Close(ctx)

		users := []models.User{}
		for cursor.Next(ctx) {
			var u models.User
			if err := cursor.Decode(&u); err != nil {
				return user, err
			}
			users = append(users, u)
		}

		if len(users) == 1 {
			return users[0], nil
		} else if len(users) > 1 {
			return user, fmt.Errorf("multiple users found with the name '%s'. Please use email or mobile number to identify the user", identifier)
		} else {
			return user, fmt.Errorf("no user found with the identifier '%s'", identifier)
		}
	}
}
