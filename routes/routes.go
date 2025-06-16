package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yashaswini7291/Inventory/controllers"
)

func UserRoutes(inRoute *gin.Engine) {
	inRoute.POST("/register", controllers.SignUp())
	inRoute.POST("/login", controllers.Login())
}
