package helper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/ttacon/libphonenumber"
)

var ErrorRecordNotFound = errors.New("record not found")
var CountryCode = "MM"

func IsValidEmail(email string) bool {
    // Basic email validation regex pattern
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    regex := regexp.MustCompile(pattern)
    return regex.MatchString(email)
}

func IsRecordValidByID(id uint, model interface{}, db *gorm.DB) bool {

    modelType := reflect.TypeOf(model).Elem() // Get the type of the element (struct)
    record := reflect.New(modelType).Interface()
    // Construct a query using the model's primary key
    query := db.Where("id = ?", id)

    // Perform the query
    if err := query.First(record).Error; err != nil {
        return false // Record with the given ID does not exist
    }

    return true
}

func ValidatePhoneNumber(phoneNumber, countryCode string) error {
    p, err := libphonenumber.Parse(phoneNumber, countryCode)
    if err != nil {
        return err // Phone number is invalid
    }

    if !libphonenumber.IsValidNumber(p) {
        return fmt.Errorf("phone number is not valid")
    }

    return nil // Phone number is valid for the specified country code
}
