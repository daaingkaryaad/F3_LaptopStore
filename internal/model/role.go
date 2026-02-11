package model

type Role struct {
	ID          int    `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
