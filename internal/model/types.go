package model

import "time"

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Permission struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

type RolePermission struct {
	RoleID       int `json:"role_id"`
	PermissionID int `json:"permission_id"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	RoleID    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Brand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Specifications struct {
	CPU              string `json:"cpu"`
	RAM              string `json:"ram"`
	Storage          string `json:"storage"`
	StorageType      string `json:"storage_type"`
	GPU              string `json:"gpu"`
	ScreenSize       string `json:"screen_size"`
	ScreenResolution string `json:"screen_resolution"`
}

type Product struct {
	ID          int            `json:"id"`
	ModelName   string         `json:"model_name"`
	BrandID     int            `json:"brand_id"`
	CategoryID  int            `json:"category_id"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	Description string         `json:"description,omitempty"`
	IsActive    bool           `json:"is_active"`
	Specs       Specifications `json:"specs"`
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
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}
