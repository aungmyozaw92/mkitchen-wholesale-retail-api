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
	Title             					string    			`gorm:"size:255;not null:unique" json:"title" validate:"required,min=3,max=200"`
	Description       					string    			`gorm:"type:text;not null" json:"description" validate:"required,min=3"`
	Price   							float64   			`gorm:"type:decimal(10,2);not null;default:0.0" json:"price" validate:"required"`
	ComparePrice   					    float64   			`gorm:"type:decimal(10,2);default:0.0" json:"compare_price"`
	Cost   					            float64   			`gorm:"type:decimal(10,2);default:0.0" json:"cost"`
	SKU                             	string    			`gorm:"size:100;not null;unique" json:"sku"  validate:"required,min=3,max=50"`
    Barcode                         	string    			`gorm:"size:100;unique" json:"barcode"  validate:"required,min=3,max=50"`
	IsQtyTracked 	 					bool 	  			`gorm:"default:false" json:"is_qty_tracked"`
	IsPhysicalProduct 	 				bool 	  			`gorm:"default:false" json:"is_physical_product"`
	IsContinueSellingWhenOutOfStock 	bool 	  			`gorm:"default:false" json:"is_continue_selling_when_out_of_stock"`
	Weight                              float64   			`gorm:"type:decimal(10,2);default:0.0" json:"weight"`
	ProductCategory   					*ProductCategory 	`gorm:"foreignKey:ProductCategoryId" json:"product_category"`
	ProductCategoryId 					uint            	`gorm:"index;not null" json:"product_category_id" validate:"required"`
	Supplier   							*Supplier 			`gorm:"foreignKey:SupplierId" json:"supplier"`
	SupplierId 							uint            	`gorm:"index;not null" json:"supplier_id" validate:"required"`
	Images      						[]Image 			`gorm:"polymorphic:Owner"`
	ProductOptions 						[]ProductOption     `json:"product_options" validate:"required,dive,required"`
	ProductVariations 					[]ProductVariation  `json:"product_variations" validate:"required,dive,required"`
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

func GetAllProducts() ([]Product, error) {

	var results []Product

	if err := DB.Find(&results).Error; err != nil {
		return results, errors.New("no products")
	}
	return results, nil
}

func GetProduct(id uint64) (Product, error) {

	var result Product

	err := DB.Preload("ProductCategory").
			Preload("Supplier").
			Preload("Images").
			Preload("ProductOptions").
			Preload("ProductVariations.Images").
			First(&result, id).Error

	if err != nil {
		return result, helper.ErrorRecordNotFound
	}

    result.Images = transformImageURLs(result.Images)
    for i := range result.ProductVariations {
        result.ProductVariations[i].Images = transformImageURLs(result.ProductVariations[i].Images)
    }

	return result, nil
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

	err := DB.Model(&Product{}).Where("sku = ? OR barcode = ? OR title = ? ", input.SKU, input.Barcode, input.Title).
				Count(&count).Error
	if err != nil {
		return &Product{}, err
	}
	if count > 0 {
		return &Product{}, errors.New("duplicate sku or barcode or product title")
	}

	images, err := uploadImages(input.Images)
	if err != nil {
		return &Product{}, err
	}

	input.Images = images

	var productOptions []ProductOption

	for _, optionRequest := range input.ProductOptions {
        productOption := ProductOption{
            OptionName:  optionRequest.OptionName,
            OptionValue: optionRequest.OptionValue,
        }
        productOptions = append(productOptions, productOption)
    }

	input.ProductOptions = productOptions

	var productVariations []ProductVariation

	for _, variation := range input.ProductVariations {

		images, err := uploadImages(variation.Images)
			if err != nil {
				return &Product{}, err
			}
        productVariation := ProductVariation{
			
            VariantName:  variation.VariantName,
            Price:  		variation.Price,
            SKU:  			variation.SKU,
            Barcode:  		variation.Barcode,
            Images: 		images,
        }
        productVariations = append(productVariations, productVariation)
    }

	input.ProductVariations = productVariations

	err = DB.Create(&input).Error
	
	if err != nil {
		return &Product{}, err
	}
	return input, nil
}

func uploadImages(images []Image) ([]Image, error) {
    var uploadedImages []Image

    for _, image := range images {
        uploadedImage, err := uploadAndAppendImage(image)
        if err != nil {
            return nil, err
        }
        uploadedImages = append(uploadedImages, uploadedImage...)
    }

    return uploadedImages, nil
}

func uploadAndAppendImage(image Image) ([]Image, error) {
	
	var images []Image

	storagePath := "products/"
	uniqueFilename := helper.GenerateUniqueFilename()
	imageFilePath := filepath.Join(storagePath, uniqueFilename)
	imageObjectUrl := storagePath + uniqueFilename

    err := utils.SaveImageToSpaces(imageObjectUrl,  image.ImageUrl)
	if err != nil {
		return images, errors.New("failed upload to cloud space")
	}

	imageObject := Image{
		ImageUrl:  imageFilePath,
	}

	images = append(images, imageObject)

	return images, nil
}

func (input *Product) DeleteProduct(id uint64) (*Product, error) {

	err := DB.Model(&Product{}).Where("id = ?", id).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

	err = DB.Delete(&input).Error
	if err != nil {
		return &Product{}, err
	}
	return input, nil
}