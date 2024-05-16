package routes

import (
	"github.com/EnesDemirtas/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItem(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderItems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:id", controllers.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems", controllers.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:id", controllers.UpdateOrderItem())
}