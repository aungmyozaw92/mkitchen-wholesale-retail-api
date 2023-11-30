package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderItem struct {
	ID                	uint      	   			`gorm:"primary_key" json:"id"`
	PurchaseOrder   	*PurchaseOrder 			`gorm:"foreignKey:PurchaseOrderId" json:"purchase_order"`
	PurchaseOrderId 	uint            		`gorm:"index;not null" json:"purchase_order_id"`
	ProductVariation   	*ProductVariation 		`gorm:"foreignKey:ProductVariationId" json:"product_variation"`
	ProductVariationId 	uint            		`gorm:"index;not null" json:"product_variation_id" validate:"required"`
	SupplierSKU         string    				`gorm:"size:255;" json:"supplier_sku"`
	ProductName         string    				`gorm:"size:255;not null" json:"product_name" validate:"required"`
	Qty   		        float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"qty" validate:"required"`
	UnitPrice   		float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"unit_price" validate:"required"`
	TaxAmount   		float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"tax_amount"`
	TaxPercent   		*float64    			`json:"tax_percent"`
	TotalAmount   		float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"total_amount"`
	ReceivedStatus      Status 					`gorm:"type:enum('pending', 'partial', 'complete');default:'pending'" json:"received_status"`
	TotalReceivedQty    float64    				`gorm:"" json:"total_received_qty"`
	TotalRemainingQty   float64    				`gorm:"" json:"total_remaining_qty"`
	CreatedAt   		time.Time				`json:"created_at"`
	UpdatedAt   		time.Time				`json:"updated_at"`
	DeletedAt        	gorm.DeletedAt   		`gorm:"index"`

}

type ReceivePurchaseOrderItem struct {
	ID                	uint      	   			`gorm:"primary_key" json:"id"`
	ReceivedQty    		float64    				`gorm:"" json:"receive_qty" validate:"required"`
	ReceivedStatus      Status 					`gorm:"type:enum('pending', 'partial', 'complete');default:'pending'" json:"received_status"`
	TotalReceivedQty    float64    				`gorm:"" json:"total_received_qty"`
	TotalRemainingQty   float64    				`gorm:"" json:"total_remaining_qty"`
}

