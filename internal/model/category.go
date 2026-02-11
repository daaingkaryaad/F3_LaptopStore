package model

type Category struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
