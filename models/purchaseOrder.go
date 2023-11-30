package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils"
	"gorm.io/gorm"
)

type Status string

const (
	Pending      Status = "pending"
	Partial 	 Status = "partial"
	Complete     Status = "complete"
)

type PurchaseOrder struct {
	ID                	uint      	   			`gorm:"primary_key" json:"id"`
	OrderNo             string    				`gorm:"index;size:255;unique" json:"order_no"`
	Supplier   			*Supplier 				`gorm:"foreignKey:SupplierId" json:"supplier"`
	SupplierId 			uint            		`gorm:"index;not null" json:"supplier_id" validate:"required"`
	TotalQty   			float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"total_qty"`
	TotalTaxAmount   	float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"total_tax_amount"`
	SubTotal   			float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"sub_total"`
	TotalAmount   		float64   				`gorm:"type:decimal(10,2);not null;default:0.0" json:"total_amount"`
	Status      		Status 					`gorm:"type:enum('pending', 'partial', 'complete');default:'pending'" json:"status"`
	ReceivedStatus      Status 					`gorm:"type:enum('pending', 'partial', 'complete');default:'pending'" json:"received_status"`
	Description       	string    				`gorm:"type:text" json:"description"`
	TotalItemCount      uint    				`gorm:"" json:"total_item_count"`
	TotalReceivedQty    float64    				`gorm:"" json:"total_received_qty"`
	TotalRemainingQty   float64    				`gorm:"" json:"total_remaining_qty"`
	PurchaseOrderItems []PurchaseOrderItem 		`json:"purchase_order_items" validate:"required,dive,required"`
	PurchaseDate		time.Time 				`gorm:"" json:"purchase_date" validate:"required"`
	ReferenceNo          string    				`gorm:"size:255;" json:"reference_no"`
	NoteToSupplier       string    				`gorm:"type:text;" json:"note_to_supplier"`
	CreatedAt   		time.Time 				`json:"created_at"`
	UpdatedAt   		time.Time 				`json:"updated_at"`
	DeletedAt        	gorm.DeletedAt   		`gorm:"index"`
}

type UpdatePurchaseOrder struct {
	SupplierId     		uint                   `json:"supplier_id" validate:"required"`
	PurchaseDate	  	 time.Time 				`gorm:"" json:"purchase_date" validate:"required"`
	Description       	 string    				`gorm:"type:text" json:"description"`
	ReferenceNo          string    				`gorm:"size:255;" json:"reference_no"`
	NoteToSupplier       string    				`gorm:"type:text;" json:"note_to_supplier"`
	AddItems     		[]PurchaseOrderItem 	`json:"add_items" validate:"required,dive,required"`
	UpdateItems  		[]PurchaseOrderItem 	`json:"update_items" validate:"required,dive,required"`
	DeleteItems     	[]uint               	`json:"delete_items"`
}

type ReceivePurchaseOrder struct {
	ReceiveItems     	[]ReceivePurchaseOrderItem 	`json:"receive_items" validate:"required,dive,required"`
}

func (p *PurchaseOrder) UnmarshalJSON(data []byte) error {
    type Alias PurchaseOrder
    aux := &struct {
        PurchaseDate string `json:"purchase_date"`
        *Alias
    }{
        Alias: (*Alias)(p),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    // Custom handling for the date format
    parsedTime, err := time.Parse("2006-01-02", aux.PurchaseDate)
    if err != nil {
        return err
    }
    p.PurchaseDate = parsedTime

    return nil
}

func (p *PurchaseOrder) BeforeSave(*gorm.DB) error {
	var lastOrder PurchaseOrder

	result := DB.Unscoped().Last(&lastOrder)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if lastOrder.ID != 0 {
		p.OrderNo = fmt.Sprintf("P%05d", lastOrder.ID+1)
	} else {
		p.OrderNo = fmt.Sprintf("P%05d", 1)
	}
	
	return nil
}

func (item *PurchaseOrderItem) CalculateTaxAndTotal() {
	// Calculate tax amount
	if item.TaxPercent != nil {
		item.TaxAmount = (item.Qty * item.UnitPrice * (*item.TaxPercent)) / 100
	} else {
		item.TaxAmount = 0
	}
	// Calculate total amount
	item.TotalAmount = item.Qty*item.UnitPrice + item.TaxAmount
}

func (po *PurchaseOrder) CalculateTotals() {

	var totalQty, subTotal, totalTaxAmount, totalAmount float64
	var totalItemCount uint

	for _, item := range po.PurchaseOrderItems {
		item.CalculateTaxAndTotal()
		totalItemCount += 1
		totalQty += item.Qty
		subTotal += item.Qty * item.UnitPrice
		totalTaxAmount += item.TaxAmount
		totalAmount += item.TotalAmount
	}

	po.TotalQty = totalQty
	po.SubTotal = subTotal
	po.TotalTaxAmount = totalTaxAmount
	po.TotalAmount = totalAmount
}

func GetAllPurchaseOrders(c *gin.Context) ([]PurchaseOrder, error) {

	var results []PurchaseOrder

	pageParam := c.Query("page")
	perPageParam := c.Query("perPage")
	sortBy := c.Query("sortBy")
	orderBy := c.Query("orderBy")
	search := c.Query("search")
	status := c.Query("status")
	supplierId := c.Query("supplier_id")

	if search != "" {
		DB = DB.Where("order_no LIKE ?", "%"+search+"%").
				Or("description LIKE ?", "%"+search+"%").
				Or("reference_no LIKE ?", "%"+search+"%").
				Or("note_to_supplier LIKE ?", "%"+search+"%")
	}
	if status != "" {
		DB = DB.Where("status", status)
	}
	if supplierId != "" {
		DB = DB.Where("supplier_id", supplierId)
	}

	if err := DB.Preload("Supplier").Find(&results).Error; err != nil {
		return results, err
	}

	if err := utils.Paginate(DB, pageParam, perPageParam, &results, sortBy, orderBy); err != nil {
		return results, errors.New("no purchase orders")
	}

	return results, nil
}

func GetPurchaseOrder(id uint64) (PurchaseOrder, error) {

	var result PurchaseOrder

	err := DB.Preload("Supplier").
			Preload("PurchaseOrderItems").
			First(&result, id).Error

	if err != nil {
		return result, helper.ErrorRecordNotFound
	}

	return result, nil
}

func (input *PurchaseOrder) CreatePurchaseOrder() (*PurchaseOrder, error) {
	
    isValidSupplierId := helper.IsRecordValidByID(input.SupplierId, &Supplier{}, DB)

	if !isValidSupplierId {
		return &PurchaseOrder{}, errors.New("invalid supplier id")
	}

	var totalItemCount uint
	var totalQty, subTotal, totalAmount, totalTaxAmount float64
	var purchaseOrderItems []PurchaseOrderItem

	// Create PurchaseOrderItems and calculate tax and total amounts
	for _, item := range input.PurchaseOrderItems {
		isValidId := helper.IsRecordValidByID(item.ProductVariationId, &ProductVariation{}, DB)

		if !isValidId {
			return &PurchaseOrder{}, errors.New("invalid product variation id")
		}

		purchaseOrderItem := PurchaseOrderItem{
			SupplierSKU: 		item.SupplierSKU,
			ProductVariationId: item.ProductVariationId,
			ProductName:        item.ProductName,
			Qty:                item.Qty,
			TotalRemainingQty:  item.Qty,
			UnitPrice:          item.UnitPrice,
			TaxPercent:         item.TaxPercent,
		}
		// Calculate tax and total amounts for the item
		purchaseOrderItem.CalculateTaxAndTotal()

		// Add the item to the PurchaseOrder
		purchaseOrderItems = append(purchaseOrderItems, purchaseOrderItem)

		totalItemCount += 1
        totalQty += purchaseOrderItem.Qty
        subTotal += purchaseOrderItem.Qty * purchaseOrderItem.UnitPrice
        totalTaxAmount += purchaseOrderItem.TaxAmount
        totalAmount += purchaseOrderItem.TotalAmount
	}

	input.PurchaseOrderItems = purchaseOrderItems
	input.TotalItemCount = totalItemCount
	input.TotalQty = totalQty
	input.SubTotal = subTotal
	input.TotalTaxAmount = totalTaxAmount
	input.TotalAmount = totalAmount

	err := DB.Create(&input).Error

	if err != nil {
		return &PurchaseOrder{}, err
	}
	return input, nil
}

func (input *UpdatePurchaseOrder) UpdatePurchaseOrder(id uint64) (*PurchaseOrder, error) {

	tx := DB.Begin()

    var existingPurchaseOrder PurchaseOrder
	if err := tx.First(&existingPurchaseOrder, id).Error; err != nil {
		return &PurchaseOrder{}, errors.New("error fetching purchase order")
	}

	// Update purchase order fields with the payload
    existingPurchaseOrder.SupplierId = input.SupplierId
    existingPurchaseOrder.PurchaseDate = input.PurchaseDate
    existingPurchaseOrder.Description = input.Description
    existingPurchaseOrder.ReferenceNo = input.ReferenceNo
    existingPurchaseOrder.NoteToSupplier = input.NoteToSupplier

    // Process add_items

    for _, addItem := range input.AddItems {
		isValidId := helper.IsRecordValidByID(addItem.ProductVariationId, &ProductVariation{}, DB)

		if !isValidId {
			return &PurchaseOrder{}, errors.New("invalid product variation id")
		}
        newItem := PurchaseOrderItem{
            ProductVariationId: addItem.ProductVariationId,
            ProductName:        addItem.ProductName,
            Qty:                addItem.Qty,
            UnitPrice:          addItem.UnitPrice,
            TaxPercent:         addItem.TaxPercent,
        }
		newItem.CalculateTaxAndTotal()
        existingPurchaseOrder.PurchaseOrderItems = append(existingPurchaseOrder.PurchaseOrderItems, newItem)
    }

    // Process update_items
	
    for _, updateItem := range input.UpdateItems {
		var existingItem PurchaseOrderItem

		if err := tx.Where("ID = ? AND purchase_order_id = ?", updateItem.ID, id).First(&existingItem).Error; err !=  nil {         
			tx.Rollback()
			return &PurchaseOrder{}, err
		}

		existingItem.ProductVariationId = updateItem.ProductVariationId
		existingItem.ProductName = updateItem.ProductName
		existingItem.Qty = updateItem.Qty
		existingItem.UnitPrice = updateItem.UnitPrice
		existingItem.TaxPercent = updateItem.TaxPercent
		
		existingItem.CalculateTaxAndTotal()
		
		if err := tx.Save(&existingItem).Error; err != nil {
			tx.Rollback()
			return &PurchaseOrder{}, err
		}
		existingPurchaseOrder.PurchaseOrderItems = append(existingPurchaseOrder.PurchaseOrderItems, updateItem)
    }

    // Process delete_items

	for _, deleteItemID := range input.DeleteItems {

		var existingItem PurchaseOrderItem

		if err := tx.Where("ID = ? AND purchase_order_id = ?", deleteItemID, id).First(&existingItem).Error; err != nil {
			tx.Rollback()
			return &PurchaseOrder{}, err
		}

		if err := tx.Delete(&existingItem).Error; err != nil {
			tx.Rollback()
			return &PurchaseOrder{}, err
		}
	}

    // Update total quantities and amounts
    existingPurchaseOrder.CalculateTotals()

    // Save the updated purchase order
    if err := tx.Save(&existingPurchaseOrder).Error; err != nil {
		tx.Rollback()
        return &PurchaseOrder{}, err
    }

	if err := tx.Commit().Error; err != nil {
        return &PurchaseOrder{}, err
    }

    return &existingPurchaseOrder, nil
}

func (input *ReceivePurchaseOrder) ReceivePurchaseOrder(id uint64) (*PurchaseOrder, error) {

	tx := DB.Begin()

    var existingPurchaseOrder PurchaseOrder
	if err := tx.First(&existingPurchaseOrder, id).Error; err != nil {
		return &PurchaseOrder{}, errors.New("error fetching purchase order")
	}

	if existingPurchaseOrder.ReceivedStatus == Complete && existingPurchaseOrder.TotalRemainingQty == 0 {
		return &PurchaseOrder{}, errors.New("this purchase order is already received")
	}

    // Process update_items
	
    for _, updateItem := range input.ReceiveItems {
		var existingItem PurchaseOrderItem

		if err := tx.Where("ID = ? AND purchase_order_id = ?", updateItem.ID, id).First(&existingItem).Error; err !=  nil {         
			tx.Rollback()
			return &PurchaseOrder{}, err
		}

		if updateItem.ReceivedQty > existingItem.TotalRemainingQty {
			tx.Rollback()
			return &PurchaseOrder{}, errors.New("please enter receive qty less than remaining qty")
		}

		existingItem.TotalReceivedQty += updateItem.ReceivedQty
		existingItem.TotalRemainingQty = existingItem.TotalRemainingQty - updateItem.ReceivedQty

		if existingItem.TotalRemainingQty > 0 {
			existingItem.ReceivedStatus = Partial
		}else{
			existingItem.ReceivedStatus = Complete
		}
		
		if err := tx.Save(&existingItem).Error; err != nil {
			tx.Rollback()
			return &PurchaseOrder{}, err
		}

		existingPurchaseOrder.TotalReceivedQty += updateItem.ReceivedQty
		existingPurchaseOrder.TotalRemainingQty = existingPurchaseOrder.TotalRemainingQty - updateItem.ReceivedQty
		
    }

	if existingPurchaseOrder.TotalRemainingQty > 0 {
		existingPurchaseOrder.ReceivedStatus = Partial
	}else{
		existingPurchaseOrder.ReceivedStatus = Complete
	}

    // Save the updated purchase order
    if err := tx.Save(&existingPurchaseOrder).Error; err != nil {
		tx.Rollback()
        return &PurchaseOrder{}, err
    }

	if err := tx.Commit().Error; err != nil {
        return &PurchaseOrder{}, err
    }

    return &existingPurchaseOrder, nil
}

func (input *PurchaseOrder) DeletePurchaseOrder(id uint64) (*PurchaseOrder, error) {
    
    tx := DB.Begin()

    if err := tx.First(input, id).Error; err != nil {
        tx.Rollback()
        return nil, helper.ErrorRecordNotFound
    }

	if err := tx.Model(&input).Association("PurchaseOrderItems").Unscoped().Clear(); err != nil {
    	tx.Rollback()
        return nil, err
    }

    if err := tx.Delete(&input).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return input, nil
}
