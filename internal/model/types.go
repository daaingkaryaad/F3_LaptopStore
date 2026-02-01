package model

import "time"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	RoleID   int    `json:"roleID"`
}

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Specifications struct {
	CPU string `json:"cpu"`
	RAM string `json:"ram"`
}

type Product struct {
	ID         int            `json:"id"`
	ModelName  string         `json:"model_name"`
	Price      float64        `json:"price"`
	BrandID    int            `json:"brand_id"`
	CategoryID int            `json:"category_id"`
	Specs      Specifications `json:"specs"`
}

type CartItem struct {
	LaptopID int `json:"laptop_id"`
	Quantity int `json:"quantity"`
}

type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type OrderItem struct {
	LaptopID int     `json:"laptop_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID     int         `json:"id"`
	UserID int         `json:"user_id"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
	Time   time.Time   `json:"time"`
}
