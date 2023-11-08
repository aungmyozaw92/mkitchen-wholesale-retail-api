package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
)

func GetAllSuppliers(context *gin.Context) {

	data, err := models.GetAllSuppliers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "success", "data": data})
}

func GetSupplier(context *gin.Context) {

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Supplier ID"})
        return
    }

	model, err := models.GetSupplier(id)
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

func CreateSupplier(context *gin.Context) {

	var input models.Supplier
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	if err := validator.New().Struct(input); err != nil {
		errorResponse := helper.ProcessValidationErrors(err)

        context.JSON(http.StatusBadRequest, gin.H{"error": errorResponse})
        return
	}

	_, err := input.CreateSupplier()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "create success"})
	
}

func UpdateSupplier(context *gin.Context) {

	var input models.Supplier

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
        return
    }

	if err := validator.New().Struct(input); err != nil {
		errorResponse := helper.ProcessValidationErrors(err)

        context.JSON(http.StatusBadRequest, gin.H{"error": errorResponse})
        return
	}

	_, err = input.UpdateSupplier(id)

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

func DeleteSupplier(context *gin.Context) {

	var input models.Supplier
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
        return
    }
	
	_, err = input.DeleteSupplier(id)
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