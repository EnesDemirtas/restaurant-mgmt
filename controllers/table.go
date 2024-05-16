package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		allTables, err := services.GetTables(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		table, err := services.GetTable(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.CreateTable(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.UpdateTable(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)

	}
}
