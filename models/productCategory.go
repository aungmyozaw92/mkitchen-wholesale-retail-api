package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
)

type ProductCategory struct {
	ID               uint     			`gorm:"primary_key" json:"id"`
	Name             string   			`gorm:"size:255;not null:unique" json:"name" binding:"required"`
	NameMM           string    			`gorm:"size:255;not null:unique" json:"name_mm" binding:"required"`
	ParentCategoryId *uint     			`gorm:"index" json:"parent_category_id"`
	SubCategories    []ProductCategory  `gorm:"foreignkey:ParentCategoryId" json:"sub_categories"`
	CreatedAt        time.Time 			`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time 			`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CategoryRelation struct {
	ID               uint    
	Name             string  
	NameMM           string   
	ParentCategoryId *uint    
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (input *ProductCategory) BeforeSave() error {
	//remove spaces
	input.Name = html.EscapeString(strings.TrimSpace(input.Name))
	input.NameMM = html.EscapeString(strings.TrimSpace(input.NameMM))

	return nil
}

func GetAllProductCategories() ([]ProductCategory, error) {

	var results []ProductCategory

	// if err := DB.Find(&results).Error; err != nil {
	if err := DB.Preload("SubCategories").Find(&results).Error; err != nil {
		return results, errors.New("no product categories")
	}
	return results, nil
}

func GetProductCategory(id uint64) (ProductCategory, error) {

	var result ProductCategory

	err := DB.Preload("SubCategories").First(&result, id).Error

	if err != nil {
		return result, helper.ErrorRecordNotFound
	}

	return result, nil
}

func (input *ProductCategory) CreateProductCategory() (*ProductCategory, error) {

	isValidId := helper.IsRecordValidByID(*input.ParentCategoryId, &ProductCategory{}, DB)

    if input.ParentCategoryId != nil && !isValidId {
        return &ProductCategory{}, errors.New("invalid product category id")
    }

	var count int64

	err := DB.Model(&ProductCategory{}).Where("name = ?", input.Name).Or("name_mm = ?", input.NameMM).Count(&count).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	if count > 0 {
		return &ProductCategory{}, errors.New("duplicate name or name_mm")
	}

	err = DB.Create(&input).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	return input, nil
}

func (input *ProductCategory) UpdateProductCategory(id uint64) (*ProductCategory, error) {

	isValidId := helper.IsRecordValidByID(*input.ParentCategoryId, &ProductCategory{}, DB)

    if input.ParentCategoryId != nil && !isValidId {
		return &ProductCategory{}, errors.New("invalid product category id")
	}

	var count int64

	err := DB.Model(&ProductCategory{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		return &ProductCategory{}, err
	}

	if count <= 0 {
		return &ProductCategory{}, helper.ErrorRecordNotFound
	}

    if err = DB.Model(&ProductCategory{}).
        Where("name = ? OR name_mm = ?", input.Name, input.NameMM).
        Not("id = ?", id).
        Count(&count).Error; err != nil {
        return nil, err
    }

    if count > 0 {
        return nil, errors.New("duplicate name or name mm")
    }

    err = DB.Model(&input).Where("id = ?", id).
        Updates(ProductCategory{Name: input.Name, NameMM: input.NameMM,ParentCategoryId: input.ParentCategoryId}).Error

    if err != nil {
        return nil, err
    }

    return input, nil
}

func (input *ProductCategory) DeleteProductCategory(id uint64) (*ProductCategory, error) {

	err := DB.Model(&ProductCategory{}).Where("id = ?", id).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

	var count int64

	err = DB.Model(&ProductCategory{}).Where("parent_category_id = ?", id).Count(&count).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	if count > 0 {
		return &ProductCategory{}, errors.New("category has sub-categories")
	}

	err = DB.Delete(&input).Error
	if err != nil {
		return &ProductCategory{}, err
	}
	return input, nil
}