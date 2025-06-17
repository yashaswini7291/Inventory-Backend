package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yashaswini7291/Inventory/routes"

	_ "github.com/yashaswini7291/Inventory/docs" // ✅ import for Swagger docs
)

// @title Inventory Management API
// @version 1.0
// @description This is a sample server for Inventory Management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@inventory.local

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	router := gin.New()
	router.Use(gin.Logger())

	// Public routes
	routes.UserRoutes(router)

	// Swagger docs
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes
	routes.ProductRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Fatal(router.Run(":" + port))
}
