package models

type Document struct {
	ID      string `json:"id", bson:"_id,omitempty"`
	UserID  string `json:"userID", bson:"userID"`
	Aadhar  string `json:"aadhar", bson:"aadhar"` // Hashed Aadhar Number
	FileURL string `json:"fileURL", bson:"fileURL"`
}
