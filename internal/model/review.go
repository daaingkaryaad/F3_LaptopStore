package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	LaptopID  string             `json:"laptop_id" bson:"laptop_id"`
	Rating    int                `json:"rating" bson:"rating"`
	Comment   string             `json:"comment" bson:"comment"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
