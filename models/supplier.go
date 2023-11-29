package models

import (
	"errors"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"gorm.io/gorm"
)

type Supplier struct {
	ID          uint      		`gorm:"primary_key" json:"id"`
	Name        string    		`gorm:"size:255;not null" json:"name" validate:"required,min=3,max=50"`
	Email       string    		`gorm:"size:255;unique" json:"email" validate:"required,email"`
	Address     string    		`gorm:"type:text;not null" json:"address" validate:"required"`
	Phone       string    		`gorm:"size:255;unique;not null" json:"phone" validate:"required,min=5,max=16"`
	Password    string    		`gorm:"size:100" json:"password"`
	CreatedAt   time.Time 		`json:"created_at"`
	UpdatedAt   time.Time 		`json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}

func (result *Supplier) PrepareGive() {

	result.Password = ""
}

func (input *Supplier) BeforeSave(*gorm.DB) error {

	//turn password into hash
	if input.Password != "" {
		hashedPassword, err := Hash(input.Password)
			if err != nil {
				return err
			}
			input.Password = string(hashedPassword)
	}

	return nil
}

func GetAllSuppliers() ([]Supplier, error) {

	var results []Supplier

	if err := DB.Find(&results).Error; err != nil {
		return results, errors.New("no Supplier")
	}

	for i, u := range results {
		u.Password = ""
		results[i] = u
	}
	return results, nil
}

func GetSupplier(id uint64) (Supplier, error) {

	var result Supplier

	err := DB.First(&result, id).Error

	if err != nil {
		return result, helper.ErrorRecordNotFound
	}

	result.PrepareGive()

	return result, nil
}

func (input *Supplier) CreateSupplier() (*Supplier, error) {

	var count int64

	if  input.Email != "" && !helper.IsValidEmail(input.Email) {
		return &Supplier{}, errors.New("invalid email address")
	}

    if err := helper.ValidatePhoneNumber(input.Phone, helper.CountryCode); err != nil {
        return &Supplier{}, errors.New("phone number validation error")
    } 

	err := DB.Model(&Supplier{}).Where("phone = ?", input.Phone).Or("email = ?", input.Email).Count(&count).Error
	if err != nil {
		return &Supplier{}, err
	}
	if count > 0 {
		return &Supplier{}, errors.New("duplicate phone or email")
	}

	err = DB.Create(&input).Error
	if err != nil {
		return &Supplier{}, err
	}
	return input, nil
}

func (input *Supplier) UpdateSupplier(id uint64) (*Supplier, error) {

	var count int64

	if  input.Email != "" && !helper.IsValidEmail(input.Email) {
		return &Supplier{}, errors.New("invalid email address")
	}

    if err := helper.ValidatePhoneNumber(input.Phone, helper.CountryCode); err != nil {
        return &Supplier{}, errors.New("phone number validation error")
    } 

	err := DB.Model(&Supplier{}).Where("id = ?", id).Count(&count).Error

	if err != nil {
		return &Supplier{}, err
	}

	if count <= 0 {
		return &Supplier{}, helper.ErrorRecordNotFound
	}

    if err = DB.Model(&Supplier{}).
        Where("email = ? OR phone = ?", input.Email, input.Phone).
        Not("id = ?", id).
        Count(&count).Error; err != nil {
        return nil, err
    }

    if count > 0 {
        return nil, errors.New("duplicate email or phone")
    }

	if input.Password != "" {
		hashedPassword, err := Hash(input.Password)
		if err != nil {
			return nil, err
		}
		input.Password = string(hashedPassword)
	}

    err = DB.Model(&input).Where("id = ?", id).
        Updates(Supplier{Name: input.Name, Email: input.Email,Phone: input.Phone,Address: input.Address,Password: input.Password}).Error

    if err != nil {
        return nil, err
    }

    return input, nil
}

func (input *Supplier) DeleteSupplier(id uint64) (*Supplier, error) {

	err := DB.Model(&Supplier{}).Where("id = ?", id).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

	err = DB.Delete(&input).Error
	if err != nil {
		return &Supplier{}, err
	}
	return input, nil
}