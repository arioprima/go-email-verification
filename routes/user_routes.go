package routes

import (
	"github.com/gin-gonic/gin"
	"golang_email_verification/controller"
)

func UserRouter(userController *controller.UserController) *gin.Engine {
	service := gin.Default()

	router := service.Group("/api/auth")

	router.POST("/login", userController.Login)
	router.POST("/register", userController.Register)
	router.POST("/verify-email", userController.VerifyEmail)

	return service
}
