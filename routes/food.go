package routes

import (
	"github.com/EnesDemirtas/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func Food(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods/", controllers.GetFoods())
	incomingRoutes.GET("/foods/:id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PATCH("/food/:id", controllers.UpdateFood())
}