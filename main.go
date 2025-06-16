package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yashaswini7291/Inventory/controllers"
	"github.com/yashaswini7291/Inventory/middleware"
	"github.com/yashaswini7291/Inventory/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println(" Server running on port", port)
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.PUT("/products/:id/quantity", controllers.UpdateProductQuantity())
	router.GET("/products", controllers.GetAllProducts())
	router.POST("/products", controllers.AddProduct())

	log.Fatal(router.Run(":" + port))
}
