package model

type OrderItem struct {
	LaptopID string  `json:"laptop_id" bson:"laptop_id"`
	Quantity int     `json:"quantity" bson:"quantity"`
	Price    float64 `json:"price" bson:"price"`
}
