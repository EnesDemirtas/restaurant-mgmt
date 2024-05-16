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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing order items",
		}
	}

	var allOrders []bson.M
	if err := result.All(ctx, &allOrders); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allOrders, nil
}

func GetOrder(c *gin.Context) (models.Order, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderId := c.Param("order_id")
	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the order item"})
		return models.Order{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while fetching the order item",
		}
	}

	return order, nil
}

func CreateOrder(c *gin.Context) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table models.Table
	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	validationErr := validate.Struct(order)

	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	if order.TableID != nil {
		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
		if err != nil {
			return nil, helpers.HttpError{
				Code:    http.StatusInternalServerError,
				Message: "table was not found",
			}
		}
	}

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	result, insertErr := orderCollection.InsertOne(ctx, order)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "order item was not created",
		}
	}

	return result, nil
}

func UpdateOrder(c *gin.Context) (*mongo.UpdateResult, error) {
	var table models.Table
	var order models.Order

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := c.BindJSON(&order); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var updateObj primitive.D

	orderId := c.Param("order_id")

	if order.TableID != nil {
		err := orderCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
		if err != nil {
			return nil, helpers.HttpError{
				Code:    http.StatusInternalServerError,
				Message: "table was not found",
			}
		}
		updateObj = append(updateObj, primitive.E{Key: "table", Value: order.TableID})
	}

	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: order.UpdatedAt})

	upsert := true
	filter := bson.M{
		"order_id": orderId,
	}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderCollection.UpdateOne(
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
			Message: "order item update failed",
		}
	}

	return result, nil
}

func OrderItemOrderCreator(ctx context.Context, order models.Order) (string, error) {
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return order.OrderID, nil
}
