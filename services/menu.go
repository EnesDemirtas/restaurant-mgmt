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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := menuCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing the menu items",
		}
	}

	var allMenus []bson.M
	if err = result.All(ctx, &allMenus); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allMenus, nil
}

func GetMenu(c *gin.Context) (models.Menu, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	menuId := c.Param("menu_id")
	var menu models.Menu

	err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
	if err != nil {
		return models.Menu{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while fetching the menu item",
		}
	}

	return menu, nil
}

func CreateMenu(c *gin.Context) (*mongo.InsertOneResult, error) {
	var menu models.Menu
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := c.BindJSON(&menu); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	validationErr := validate.Struct(menu)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.ID = primitive.NewObjectID()
	menu.MenuID = menu.ID.Hex()

	result, insertErr := menuCollection.InsertOne(ctx, menu)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Menu item was not created",
		}
	}

	return result, nil
}

func UpdateMenu(c *gin.Context) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu

	if err := c.BindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	menuId := c.Param("menu_id")
	filter := bson.M{"menu_id": menuId}

	var updateObj primitive.D

	if menu.StartDate == nil || menu.EndDate == nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Please type the start and end date",
		}
	}

	if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "kindly retype the time",
		}
	}

	updateObj = append(updateObj, primitive.E{Key: "start_date", Value: menu.StartDate})
	updateObj = append(updateObj, primitive.E{Key: "end_date", Value: menu.EndDate})

	if menu.Name != "" {
		updateObj = append(updateObj, primitive.E{Key: "name", Value: menu.Name})
	}

	if menu.Category != "" {
		updateObj = append(updateObj, primitive.E{Key: "category", Value: menu.Category})
	}

	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: menu.UpdatedAt})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Menu update failed",
		}
	}

	return result, nil
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.Before(check) && end.After(start)
}
