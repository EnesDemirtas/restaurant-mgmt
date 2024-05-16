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
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func GetUsers(c *gin.Context) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * recordPerPage
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			},
		},
	}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, projectStage,
	})

	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	var allUsers []bson.M
	if err := result.All(ctx, &allUsers); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allUsers, nil
}

func GetUser(c *gin.Context) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userId := c.Param("user_id")

	var user models.User

	err := userCollection.FindOne(ctx, bson.M{
		"user_id": userId,
	}).Decode(&user)

	if err != nil {
		return models.User{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return user, nil
}

func SignUp(c *gin.Context) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: validationErr.Error(),
		}
	}

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while checking for the email",
		}
	}

	if emailCount > 0 {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "this email already exists",
		}
	}

	password := helpers.HashPassword(*user.Password)
	user.Password = &password

	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while checking for the phone",
		}
	}

	if phoneCount > 0 {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "this phone number already exists",
		}
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()

	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
	user.Token = &token
	user.RefreshToken = &refreshToken

	insertedID, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: insertErr.Error(),
		}
	}

	return insertedID.InsertedID, nil
}

func Login(c *gin.Context) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		return models.User{}, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		return models.User{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "user not found",
		}
	}

	passwordIsValid, msg := helpers.VerifyPassword(*user.Password, *foundUser.Password)
	if !passwordIsValid {
		return models.User{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(
		*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)

	helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID)

	return foundUser, nil
}
