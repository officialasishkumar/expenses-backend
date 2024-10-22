// main.go
package main

import (
	"expenses-backend/db"
	"expenses-backend/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB
	mongoURI := "mongodb://localhost:27017" // Update as per your setup
	db.InitMongoDB(mongoURI)
	defer db.CloseMongoDB()

	// Initialize Gin router
	router := gin.Default()

	// User routes
	router.POST("/users", handlers.CreateUser)
	router.GET("/users", handlers.GetUser) // Use query parameter 'identifier'

	// Expense routes
	router.POST("/expenses", handlers.AddExpense)
	router.GET("/expenses/user", handlers.GetUserExpenses) // Use query parameter 'identifier'
	router.GET("/expenses", handlers.GetOverallExpenses)

	// Balance Sheet
	router.GET("/balancesheet/download", handlers.DownloadBalanceSheet)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
