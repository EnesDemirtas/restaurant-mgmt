package services

import (
	"context"
	"log"
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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := tableCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing tables",
		}
	}

	var allTables []bson.M

	if err := result.All(ctx, &allTables); err != nil {
		log.Fatal(err)
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allTables, nil

}

func GetTable(c *gin.Context) (models.Table, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	tableId := c.Param("table_id")
	var table models.Table

	err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
	if err != nil {
		return models.Table{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while fetching the table item",
		}
	}

	return table, nil
}

func CreateTable(c *gin.Context) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table models.Table

	if err := c.BindJSON(&table); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	validationErr := validate.Struct(table)

	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
	table.TableID = table.ID.Hex()

	result, insertErr := tableCollection.InsertOne(ctx, table)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Table item was not created",
		}
	}

	return result, nil
}

func UpdateTable(c *gin.Context) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var table models.Table

	tableId := c.Param("table_id")

	if err := c.BindJSON(&table); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	filter := bson.M{"table_id": tableId}

	var updateObj primitive.D

	if table.NumberOfGuests != nil {
		updateObj = append(updateObj, primitive.E{Key: "number_of_guests", Value: table.NumberOfGuests})
	}

	if table.TableNumber != nil {
		updateObj = append(updateObj, primitive.E{Key: "table_number", Value: table.TableNumber})
	}

	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: table.UpdatedAt})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := tableCollection.UpdateOne(
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
			Message: "Table update failed",
		}
	}

	return result, nil
}
