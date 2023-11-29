package utils

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// Paginate retrieves a specified page of results from the database
func Paginate(db *gorm.DB, pageParam, perPageParam string, results interface{}, sortBy, orderBy string) error {
// Default values
	page := 1
	perPage := 25

	// Parse page and perPage from parameters
	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil {
			page = p
		}
	}

	if perPageParam != "" {
		pp, err := strconv.Atoi(perPageParam)
		if err == nil {
			perPage = pp
		}
	}

	offset := (page - 1) * perPage
	

	orderClause := orderClause(sortBy, orderBy)

	return db.Order(orderClause).Limit(perPage).Offset(offset).Find(results).Error
}

func orderClause(sortBy, orderBy string) string {

    validSortFields := []string{"id", "created_at"}
    validOrders := []string{"asc", "desc"}

    isValidSortField := false
    for _, field := range validSortFields {
        if sortBy == field {
            isValidSortField = true
            break
        }
    }

    if sortBy != "" && !isValidSortField {
        return ""
    }

    isValidOrder := false
    for _, order := range validOrders {
        if orderBy == order {
            isValidOrder = true
            break
        }
    }

    if orderBy != "" && !isValidOrder {
        // Return an empty string for invalid orderBy
        return ""
    }

    var orderClause string
    if sortBy != "" && orderBy != "" {
        orderClause = fmt.Sprintf("%s %s", sortBy, orderBy)
    } else {
        orderClause = "created_at desc"
    }

    return orderClause
}