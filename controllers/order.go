package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		allOrders, err := services.GetOrders(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		order, err := services.GetOrder(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.CreateOrder(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.UpdateOrder(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
