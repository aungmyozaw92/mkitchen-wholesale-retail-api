package models

import "time"

type ProductOption struct {
	ID        		uint   		`gorm:"primary_key" json:"id"`
    ProductId 		uint   		`gorm:"index;not null" json:"product_id"`
    OptionName      string 		`gorm:"size:255;not null" json:"option_name" validate:"required,min=3,max=30"`
    OptionValue     string 		`gorm:"size:255;not null" json:"option_value" validate:"required,min=3,max=30"`
	CreatedAt       time.Time  	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
