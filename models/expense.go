// models/expense.go
package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Expense struct {
    ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    Description  string                 `bson:"description" json:"description" validate:"required"`
    Amount       float64                `bson:"amount" json:"amount" validate:"required,gt=0"`
    CreatedBy    primitive.ObjectID     `bson:"created_by" json:"created_by" validate:"required"`
    SplitType    string                 `bson:"split_type" json:"split_type" validate:"required,oneof=Equal Exact Percentage"`
    Participants []primitive.ObjectID   `bson:"participants" json:"participants" validate:"required,min=1"`
    SplitDetails map[string]interface{} `bson:"split_details,omitempty" json:"split_details,omitempty"`
    CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
}
