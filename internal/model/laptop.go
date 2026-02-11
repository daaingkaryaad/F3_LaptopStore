package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Laptop struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ModelName   string             `json:"model_name" bson:"model_name"`
	BrandID     string             `json:"brand_id" bson:"brand_id"`
	CategoryID  string             `json:"category_id" bson:"category_id"`
	Price       float64            `json:"price" bson:"price"`
	Stock       int                `json:"stock" bson:"stock"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	Specs       LaptopSpec         `json:"specs" bson:"specs"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
