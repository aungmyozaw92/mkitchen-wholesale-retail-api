package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
)

// func GetAllProducts(context *gin.Context) {

// 	users, err := models.GetAllProducts()
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
// }

// func GetProduct(context *gin.Context) {

// 	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
//         return
//     }

// 	model, err := models.GetProduct(id)
// 	if err != nil {
// 		if err == helper.ErrorRecordNotFound {
//             context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//             return
//         }
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"message": "success", "data": model})
// }

func CreateProduct(context *gin.Context) {

	var input models.Product
	
	if err := context.ShouldBindJSON(&input); err != nil {
		fmt.Print(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	_, err := input.CreateProduct()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "create success"})
	
}

// func UpdateProduct(context *gin.Context) {

// 	var input models.Product
// 	if err := context.ShouldBindJSON(&input); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
//         return
//     }

// 	_, err = input.UpdateProduct(id)
// 	if err != nil {
// 		if err == helper.ErrorRecordNotFound {
//             context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//             return
//         }
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"message": "update success"})
// }

// func DeleteProduct(context *gin.Context) {

// 	var input models.Product
// 	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
//         return
//     }
	
// 	_, err = input.DeleteProduct(id)
// 	if err != nil {
// 		if err == helper.ErrorRecordNotFound {
//             context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//             return
//         }
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"message": "delete success"})
// }