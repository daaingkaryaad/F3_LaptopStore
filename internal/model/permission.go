package model

type Permission struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Code        string `json:"code" bson:"code"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
