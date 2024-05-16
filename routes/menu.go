package routes

import (
	"github.com/EnesDemirtas/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func Menu(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus", controllers.GetMenus())
	incomingRoutes.GET("/menus/:id", controllers.GetMenu())
	incomingRoutes.POST("/menus", controllers.CreateMenu())
	incomingRoutes.PATCH("/menus/:id", controllers.UpdateMenu())
}