package models

type Image struct {
	ID         uint   
    ImageUrl   string `json:"image_url"`
    OwnerType  string 
    OwnerID    uint  
}

type ImageTest struct {
	ID        uint     `gorm:"primary_key" json:"id"`
    OwnerID   uint
    OwnerType string
    URL      string // Base64-encoded image data
}