package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Token     string             `json:"token" bson:"token"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Role      string             `json:"role" bson:"role"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
