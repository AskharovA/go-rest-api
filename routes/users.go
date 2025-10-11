package routes

import (
	"AskharovA/go-rest-api/models"
	"AskharovA/go-rest-api/services"
	"AskharovA/go-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	UserService services.UserService
}

func (api *UserAPI) signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = api.UserService.CreateUser(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save the user."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

func (api *UserAPI) login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = api.UserService.ValidateCredentials(&user)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
