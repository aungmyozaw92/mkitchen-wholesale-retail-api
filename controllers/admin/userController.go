package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
)

func GetAllUsers(context *gin.Context) {

	users, err := models.GetAllUsers()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "success", "data": users})
}

func GetUser(context *gin.Context) {

	// id := context.Param("id")
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

	model, err := models.GetUser(id)
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

func CreateUser(context *gin.Context) {

	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	_, err := input.CreateUser()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "create success"})
	
}

func UpdateUser(context *gin.Context) {

	var input models.User
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

	_, err = input.UpdateUser(id)
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

func DeleteUser(context *gin.Context) {

	var input models.User
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
	
	_, err = input.DeleteUser(id)
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



