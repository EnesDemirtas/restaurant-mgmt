package controllers

import (
	"net/http"

	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/services"
	"github.com/gin-gonic/gin"
)

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		allInvoices, err := services.GetInvoices(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceView, err := services.GetInvoice(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.CreateInvoice(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := services.UpdateInvoice(c)
		if err != nil {
			c.JSON(err.(helpers.HttpError).GetFields())
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
