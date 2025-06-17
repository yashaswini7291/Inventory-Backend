package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yashaswini7291/Inventory/controllers"
	"github.com/yashaswini7291/Inventory/middleware"
)

func UserRoutes(inRoute *gin.Engine) {
	inRoute.POST("/register", controllers.SignUp())
	inRoute.POST("/login", controllers.Login())
}

func ProductRoutes(router *gin.Engine) {
	protected := router.Group("/products")
	protected.Use(middleware.Authentication())

	{
		protected.PUT("/:id/quantity", controllers.UpdateProductQuantity())
		protected.GET("", controllers.GetAllProducts())
		protected.POST("", controllers.AddProduct())
	}
}
