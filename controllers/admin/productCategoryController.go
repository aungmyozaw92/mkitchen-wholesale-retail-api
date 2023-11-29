package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
)

func GetAllProductCategories(context *gin.Context) {

	users, err := models.GetAllProductCategories(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	context.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

func GetProductCategory(context *gin.Context) {

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProductCategory ID"})
        return
    }

	model, err := models.GetProductCategory(id)
	if err != nil {
		if err == helper.ErrorRecordNotFound {
            context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "success", "data": model})
}

func CreateProductCategory(context *gin.Context) {

	var input models.ProductCategory

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	if err := validator.New().Struct(input); err != nil {
		errorResponse := helper.ProcessValidationErrors(err)

        context.JSON(http.StatusBadRequest, gin.H{"error": errorResponse})
        return
	}

	_, err := input.CreateProductCategory()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "create success"})
	
}

func UpdateProductCategory(context *gin.Context) {

	var input models.ProductCategory
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProductCategory ID"})
        return
    }

	if err := validator.New().Struct(input); err != nil {
		errorResponse := helper.ProcessValidationErrors(err)

        context.JSON(http.StatusBadRequest, gin.H{"error": errorResponse})
        return
	}

	_, err = input.UpdateProductCategory(id)
	if err != nil {
		if err == helper.ErrorRecordNotFound {
            context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "update success"})
}

func DeleteProductCategory(context *gin.Context) {

	var input models.ProductCategory
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProductCategory ID"})
        return
    }
	
	_, err = input.DeleteProductCategory(id)
	if err != nil {
		if err == helper.ErrorRecordNotFound {
            context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	context.JSON(http.StatusOK, gin.H{"message": "delete success"})
}