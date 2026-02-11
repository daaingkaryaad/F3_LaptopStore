package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Items     []CartItem         `json:"items" bson:"items"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
