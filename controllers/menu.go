package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		allMenus, err := services.GetMenus(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		menu, err := services.GetMenu(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.CreateMenu(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.UpdateMenu(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, result)
	}

}
