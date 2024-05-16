package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		allOrderItems, err := services.GetOrderItems(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems, err := services.ItemsByOrder(orderId)

		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderItem, err := services.GetOrderItem(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		insertedOrderItems, err := services.CreateOrderItem(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.UpdateOrderItem(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
