package main

import (
	"os"

	"github.com/EnesDemirtas/restaurant-management/database"
	"github.com/EnesDemirtas/restaurant-management/middlewares"
	"github.com/EnesDemirtas/restaurant-management/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.User(router)
	router.Use(middlewares.Authentication())

	routes.Food(router)
	routes.Menu(router)
	routes.Table(router)
	routes.Order(router)
	routes.OrderItem(router)
	routes.Invoice(router)

	router.Run(":" + port)
}