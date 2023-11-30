package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
)

func GetAllPurchaseOrders(context *gin.Context) {

	users, err := models.GetAllPurchaseOrders(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

func GetPurchaseOrder(context *gin.Context) {

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PurchaseOrder ID"})
        return
    }

	model, err := models.GetPurchaseOrder(id)
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

func CreatePurchaseOrder(context *gin.Context) {

	var input models.PurchaseOrder
	
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	if err := validator.New().Struct(input); err != nil {
		errorResponse := helper.ProcessValidationErrors(err)

        context.JSON(http.StatusBadRequest, gin.H{"error": errorResponse})
        return
	}

	_, err := input.CreatePurchaseOrder()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "create success"})
	
}

func UpdatePurchaseOrder(context *gin.Context) {

	var input models.UpdatePurchaseOrder
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PurchaseOrder ID"})
        return
    }

	_, err = input.UpdatePurchaseOrder(id)
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

func ReceivePurchaseOrder(context *gin.Context) {

	var input models.ReceivePurchaseOrder
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PurchaseOrder ID"})
        return
    }

	_, err = input.ReceivePurchaseOrder(id)
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

func DeletePurchaseOrder(context *gin.Context) {

	var input models.PurchaseOrder
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PurchaseOrder ID"})
        return
    }
	
	_, err = input.DeletePurchaseOrder(id)
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