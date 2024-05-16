package services

import (
	"context"
	"net/http"
	"time"

	"github.com/EnesDemirtas/restaurant-management/database"
	"github.com/EnesDemirtas/restaurant-management/helpers"
	"github.com/EnesDemirtas/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct {
	InvoiceID      string
	PaymentMethod  string
	OrderID        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing all invoices",
		}
	}

	var allInvoices []bson.M
	if err = result.All(ctx, &allInvoices); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allInvoices, nil
}

func GetInvoice(c *gin.Context) (InvoiceViewFormat, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	invoiceId := c.Param("invoice_id")

	var invoice models.Invoice

	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
	if err != nil {
		return InvoiceViewFormat{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing invoice item",
		}
	}

	var invoiceView InvoiceViewFormat

	allOrderItems, err := ItemsByOrder(invoice.OrderID)
	if err != nil {
		return InvoiceViewFormat{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	invoiceView.OrderID = invoice.OrderID
	invoiceView.PaymentDueDate = invoice.PaymentDueDate

	invoiceView.PaymentMethod = "null"
	if invoice.PaymentMethod != nil {
		invoiceView.PaymentMethod = *invoice.PaymentMethod
	}

	invoiceView.InvoiceID = invoice.InvoiceID
	invoiceView.PaymentStatus = invoice.PaymentStatus
	invoiceView.PaymentDue = allOrderItems[0]["payment_due"]
	invoiceView.TableNumber = allOrderItems[0]["table_number"]
	invoiceView.OrderDetails = allOrderItems[0]["order_items"]

	return invoiceView, nil
}

func CreateInvoice(c *gin.Context) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice models.Invoice

	if err := c.BindJSON(&invoice); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order)
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Order was not found",
		}
	}

	status := "PENDING"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &status
	}

	invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.InvoiceID = invoice.ID.Hex()

	validationErr := validate.Struct(invoice)
	if validationErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "invoice item was not created",
		}
	}

	return result, nil
}

func UpdateInvoice(c *gin.Context) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var invoice models.Invoice

	invoiceId := c.Param("invoice_id")

	if err := c.BindJSON(&invoice); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	filter := bson.M{"invoice_id": invoiceId}

	var updateObj primitive.D

	if invoice.PaymentMethod != nil {
		updateObj = append(updateObj, primitive.E{Key: "payment_method", Value: invoice.PaymentMethod})
	}

	if invoice.PaymentStatus != nil {
		updateObj = append(updateObj, primitive.E{Key: "payment_status", Value: invoice.PaymentStatus})
	}

	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: invoice.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "PENDING"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &status
	}

	result, err := invoiceCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "invoice item update failed",
		}
	}

	return result, nil
}
