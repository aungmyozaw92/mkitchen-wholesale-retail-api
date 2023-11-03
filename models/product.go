package models

import (
	"errors"

	"html"
	"path/filepath"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils"
	"gorm.io/gorm"
)

type Product struct {
	ID                					uint      			`gorm:"primary_key" json:"id"`
	Title             					string    			`gorm:"size:255;not null:unique" json:"title" binding:"required"`
	Description       					string    			`gorm:"type:text;not null" json:"description" binding:"required"`
	Price   							float64   			`gorm:"type:decimal(10,2);not null;default:0.0" json:"price" binding:"required"`
	ComparePrice   					    float64   			`gorm:"type:decimal(10,2);default:0.0" json:"compare_price"`
	Cost   					            float64   			`gorm:"type:decimal(10,2);default:0.0" json:"cost"`
	SKU                             	string    			`gorm:"size:100;not null;unique" json:"sku" binding:"required"`
    Barcode                         	string    			`gorm:"size:100;unique" json:"barcode"`
	IsQtyTracked 	 					bool 	  			`gorm:"default:false" json:"is_qty_tracked"`
	IsPhysicalProduct 	 				bool 	  			`gorm:"default:false" json:"is_physical_product"`
	IsContinueSellingWhenOutOfStock 	bool 	  			`gorm:"default:false" json:"is_continue_selling_when_out_of_stock"`
	Weight                              float64   			`gorm:"" json:"weight"`
	ProductCategoryId 					uint     			`gorm:"index;not null" json:"product_category_id"`
	CategoryRelation                    CategoryRelation     `gorm:"foreignKey:ProductCategoryId" json:"product_category"`
	SupplierId 							uint     			`gorm:"index;not null" json:"supplier_id"`
	SupplierRelation    				SupplierRelation   	`gorm:"foreignkey:SupplierId" json:"supplier"`
	Images      						[]Image 			`gorm:"polymorphic:Owner"`
	CreatedAt        					time.Time  			`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        					time.Time 			`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        					gorm.DeletedAt   	`gorm:"index"`

}

func (input *Product) BeforeSave() error {
	//remove spaces
	input.Title = html.EscapeString(strings.TrimSpace(input.Title))
	input.SKU = html.EscapeString(strings.TrimSpace(input.SKU))
	input.Barcode = html.EscapeString(strings.TrimSpace(input.Barcode))

	return nil
}


func (input *Product) CreateProduct() (*Product, error) {

	isValidId := helper.IsRecordValidByID(input.ProductCategoryId, &ProductCategory{}, DB)

	if !isValidId {
		return &Product{}, errors.New("invalid product category id")
	}

	isValidSupplierId := helper.IsRecordValidByID(input.SupplierId, &Supplier{}, DB)

	if !isValidSupplierId {
		return &Product{}, errors.New("invalid supplier id")
	}

	var count int64

	err := DB.Model(&Product{}).Where("sku = ?", input.SKU).Or("barcode = ?", input.Barcode).Count(&count).Error
	if err != nil {
		return &Product{}, err
	}
	if count > 0 {
		return &Product{}, errors.New("duplicate sku or barcode")
	}

	var images []Image
    for _, image := range input.Images {

        uploadedImages, err := uploadAndAppendImage(input, image)
        if err != nil {
            return &Product{}, err
        }
        images = append(images, uploadedImages...)
    }

	product := Product{
        Title:   input.Title,
        Description: input.Description,
        Price: input.Price,
        SKU: input.SKU,
        Barcode: input.Barcode,
        ProductCategoryId: input.ProductCategoryId,
        SupplierId: input.SupplierId,
        // Other fields
        Images: images,
    }

	err = DB.Create(&product).Error
	
	if err != nil {
		return &Product{}, err
	}
	return input, nil
}

func uploadAndAppendImage(input *Product, image Image) ([]Image, error) {
	
	var images []Image

	storagePath := ""
	uniqueFilename := helper.GenerateUniqueFilename()
	imageFilePath := filepath.Join(storagePath, uniqueFilename)

    err := utils.SaveImageToSpaces(uniqueFilename,  image.ImageUrl)
	if err != nil {
		return images, errors.New("failed upload to space")
	}

	imageObject := Image{
		ImageUrl:  imageFilePath,
		OwnerType: "Product",
		OwnerID:   input.ID, 
	}

	images = append(images, imageObject)

	return images, nil
}