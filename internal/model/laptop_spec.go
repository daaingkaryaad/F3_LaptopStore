package model

type LaptopSpec struct {
	CPU              string `json:"cpu" bson:"cpu"`
	RAM              string `json:"ram" bson:"ram"`
	Storage          string `json:"storage" bson:"storage"`
	StorageType      string `json:"storage_type" bson:"storage_type"`
	GPU              string `json:"gpu" bson:"gpu"`
	ScreenSize       string `json:"screen_size" bson:"screen_size"`
	ScreenResolution string `json:"screen_resolution" bson:"screen_resolution"`
}
