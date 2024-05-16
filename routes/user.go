package routes

import (
	"github.com/EnesDemirtas/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func User(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:id", controllers.GetUser())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
}