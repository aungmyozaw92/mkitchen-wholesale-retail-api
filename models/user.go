package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Username    string    `gorm:"size:255;not null;unique" json:"username" binding:"required"`
	Name        string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Email       string    `gorm:"size:255;unique" json:"email"`
	Password    string    `gorm:"size:100;not null" json:"password"`
	IsActive    *bool     `gorm:"not null" json:"is_active"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (result *User) PrepareGive() {

	result.Password = ""
}

func Hash(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(password, hashedPassword string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (input *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := Hash(input.Password)
	if err != nil {
		return err
	}
	input.Password = string(hashedPassword)

	//remove spaces in username
	input.Username = html.EscapeString(strings.TrimSpace(input.Username))

	return nil
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(u.ID, u.Username)

	if err != nil {
		return "", err
	}

	return token, nil

}

func GetAllUsers() ([]User, error) {

	var results []User

	if err := DB.Find(&results).Error; err != nil {
		return results, errors.New("no user")
	}

	for i, u := range results {
		u.Password = ""
		results[i] = u
	}
	return results, nil
}

func (input *User) CreateUser() (*User, error) {

	var count int64

	if  input.Email != "" && !helper.IsValidEmail(input.Email) {
		return &User{}, errors.New("invalid email address")
	}

	err := DB.Model(&User{}).Where("username = ?", input.Username).Or("email = ?", input.Email).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count > 0 {
		return &User{}, errors.New("duplicate username or email")
	}

	err = DB.Create(&input).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}

func GetUser(id uint64) (User, error) {

	var result User

	err := DB.First(&result, id).Error

	if err != nil {
		return result, helper.ErrorRecordNotFound
	}

	result.PrepareGive()

	return result, nil
}

func GetUserByID(id uint) (User, error) {

	var result User

	if err := DB.First(&result, id).Error; err != nil {
		return result, helper.ErrorRecordNotFound
	}

	result.PrepareGive()

	return result, nil
}

func (input *User) UpdateUser(id uint64) (*User, error) {

	var count int64

	err := DB.Model(&User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count <= 0 {
        return nil, helper.ErrorRecordNotFound
    }

	if err = DB.Model(&User{}).
        Where("username = ? OR email = ?", input.Username, input.Email).
        Not("id = ?", id).
        Count(&count).Error; err != nil {
        return nil, err
    }
	if count > 0 {
		return &User{}, errors.New("duplicate email or username")
	}

	err = DB.Model(&input).Updates(User{Name: input.Name,Email: input.Email,Username: input.Username, IsActive: input.IsActive}).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}

func (input *User) DeleteUser(id uint64) (*User, error) {

	err := DB.Model(&User{}).Where("id = ?", id).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

	err = DB.Delete(&input).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}

func (input *User) ChangeUserPassword() (*User, error) {

	//turn password into hash
	hashedPassword, err := Hash(input.Password)
	if err != nil {
		return &User{}, err
	}
	input.Password = string(hashedPassword)

	err = DB.Model(&User{}).Where("id = ?", input.ID).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

	err = DB.Model(&input).Updates(User{Password: input.Password}).Error
	if err != nil {
		return &User{}, err
	}
	return input, nil
}
