package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/models"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils/token"
)

func CurrentUser(context *gin.Context) {

	user_id, err := token.ExtractTokenID(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model, err := models.GetUserByID(user_id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "success", "data": model})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(context *gin.Context) {

	var input LoginInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.User{}

	model.Username = input.Username
	model.Password = input.Password

	token, err := models.LoginCheck(model.Username, model.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}