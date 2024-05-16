package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		allUsers, err := services.GetUsers(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := services.GetUser(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		insertedID, err := services.SignUp(c)
		if insertedID == nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, insertedID)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		foundUser, err := services.Login(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
