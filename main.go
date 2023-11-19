package main

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"golang_email_verification/controller"
	"golang_email_verification/initializers"
	"golang_email_verification/repository"
	"golang_email_verification/routes"
	"golang_email_verification/service"
	"log"
)

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Println("🚀 Could not load environment variables", err)
	}

	db, err := initializers.ConnectDB(&config)

	validate := validator.New()
	userRepository := repository.NewUserRepositoryImpl(db)
	userService := service.NewUserServiceImpl(userRepository, db, validate)
	userController := controller.NewUserController(userService)

	router := routes.UserRouter(userController)
	err = router.Run(":8080")
	if err != nil {
		log.Println("🚀 Could not run server", err)
	}
}
