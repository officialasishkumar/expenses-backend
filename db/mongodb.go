// db/mongodb.go
package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson" // Added bson import
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client      *mongo.Client
	UsersCol    *mongo.Collection
	ExpensesCol *mongo.Collection
)

func InitMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	Client = client
	db := client.Database("expenses_db")
	UsersCol = db.Collection("users")
	ExpensesCol = db.Collection("expenses")

	// Create indexes
	createIndexes()
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unique index on email
	_, err := UsersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create index on email: %v", err)
	}

	// Unique index on mobile_number
	_, err = UsersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"mobile_number": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Failed to create index on mobile_number: %v", err)
	}
}

func CloseMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Disconnect(ctx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}
}
