package services

import (
	"context"
	"net/http"
	"strconv"
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

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil {
		startIndex = (page - 1) * recordPerPage
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{
		{
			Key: "$group", Value: bson.D{
				{
					Key: "_id", Value: bson.D{
						{Key: "_id", Value: "null"},
					},
				},
				{
					Key: "total_count", Value: bson.D{
						{Key: "$sum", Value: 1},
					},
				},
				{
					Key: "data", Value: bson.D{
						{Key: "$push", Value: "$$ROOT"},
					},
				},
			},
		},
	}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{
					"$data", startIndex, recordPerPage,
				}}}},
			},
		},
	}

	result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: groupStage}},
		{{Key: "$project", Value: projectStage}}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "	error occured while listing food items"})
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing food items",
		}
	}

	var allFoods []bson.M
	if err := result.All(ctx, &allFoods); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allFoods, nil
}

func GetFood(c *gin.Context) (models.Food, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	foodId := c.Param("food_id")
	var food models.Food

	err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		return models.Food{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while fetching the food item",
		}
	}

	return food, nil
}

func CreateFood(c *gin.Context) (*mongo.InsertOneResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu
	var food models.Food

	if err := c.BindJSON(&food); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	validationErr := validate.Struct(food)
	if validationErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuID}).Decode(&menu)
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "menu was not found",
		}
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.FoodID = food.ID.Hex()
	var num = helpers.ToFixed(*food.Price, 2)
	food.Price = &num

	result, insertErr := foodCollection.InsertOne(ctx, food)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Food item was not created",
		}
	}

	return result, nil
}

func UpdateFood(c *gin.Context) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var menu models.Menu
	var food models.Food

	if err := c.BindJSON(&food); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	foodId := c.Param("food_id")

	var updateObj primitive.D

	if food.Name != nil {
		updateObj = append(updateObj, primitive.E{Key: "name", Value: food.Name})
	}

	if food.Price != nil {
		updateObj = append(updateObj, primitive.E{Key: "price", Value: food.Price})
	}

	if food.FoodImage != nil {
		updateObj = append(updateObj, primitive.E{Key: "food_image", Value: food.FoodImage})
	}

	if food.MenuID != nil {
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuID}).Decode(&menu)
		if err != nil {
			return nil, helpers.HttpError{
				Code:    http.StatusInternalServerError,
				Message: "menu was not found",
			}
		}
		updateObj = append(updateObj, primitive.E{Key: "menu", Value: food.Price})
	}

	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: food.UpdatedAt})

	upsert := true
	filter := bson.M{"food_id": foodId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := foodCollection.UpdateOne(
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
			Message: "food item update failed",
		}
	}

	return result, nil
}
