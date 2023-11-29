package models

import "time"

type ProductVariation struct {
	ID        		uint        `gorm:"primary_key" json:"id"`
    ProductId 		uint        `gorm:"index;not null" json:"product_id"`
    VariantName     string      `gorm:"size:255;not null" json:"variant_name" validate:"required,min=3,max=30"`
    Price   		float64   	`gorm:"type:decimal(10,2);not null;default:0.0" json:"price" validate:"required"`
    SKU             string    	`gorm:"size:100;not null;unique" json:"sku"  validate:"required,min=3,max=50"`
    Barcode         string    	`gorm:"size:100;unique" json:"barcode"  validate:"required,min=3,max=50"`
    Images      	[]Image 	`gorm:"polymorphic:Owner"`
    IsDelete 		bool 		`json:"is_delete"`
    CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
