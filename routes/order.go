package routes

import (
	"github.com/EnesDemirtas/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func Order(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:id", controllers.UpdateOrder())
}