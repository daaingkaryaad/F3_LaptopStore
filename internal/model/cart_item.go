package model

type CartItem struct {
	LaptopID string `json:"laptop_id" bson:"laptop_id"`
	Quantity int    `json:"quantity" bson:"quantity"`
}
