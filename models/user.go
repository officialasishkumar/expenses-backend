// models/user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name" validate:"required"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	MobileNumber string             `bson:"mobile_number" json:"mobile_number" validate:"required"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
