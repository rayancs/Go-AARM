package main

import (
	"app/configs"
	"app/controller"
	handler "app/handlers"
	repo "app/repos"
	"app/services"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	// test.Test()
	db := repo.NewMongoDB(configs.GetMongoDBURI(), "dev-shared")
	// all dependencies
	authRepo := repo.NewMongoUser(db)
	// this is where the power of di works , the base layer has to be hardcoded ( technically) but layers after are to be used via di

	authService := services.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	authController := controller.NewAuthController(authHandler, router)
	authController.AuthEndpoints()
	controller.BasePage(router)
	router.Run(":9090")
}
