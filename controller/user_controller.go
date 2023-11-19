package controller

import (
	"github.com/gin-gonic/gin"
	"golang_email_verification/models"
	"golang_email_verification/service"
	"log"
	"net/http"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (controller *UserController) Login(ctx *gin.Context) {
	loginRequest := models.LoginInput{}
	err := ctx.ShouldBindJSON(&loginRequest)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := controller.UserService.Login(ctx, loginRequest)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Internal Server Error",
		})
	} else {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "OK",
			"Data":    loginResponse,
		})
	}
}

func (controller *UserController) Register(ctx *gin.Context) {
	registerRequest := models.RegisterInput{}
	err := ctx.ShouldBindJSON(&registerRequest)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	registerResponse, err := controller.UserService.Register(ctx, registerRequest)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Internal Server Error",
		})
	} else {
		ctx.IndentedJSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "OK",
			"Data":    registerResponse,
		})
	}
}

func (controller *UserController) VerifyEmail(ctx *gin.Context) {
	verifyRequest := models.VerifyInput{}
	if err := ctx.ShouldBindJSON(&verifyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verifyResponse, err := controller.UserService.VerifyEmail(ctx, verifyRequest)
	if err != nil {
		log.Printf("Error verifying email: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OK",
		"data":    verifyResponse,
	})
}
