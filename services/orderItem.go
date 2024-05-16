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

type OrderItemPack struct {
	TableID    *string
	OrderItems []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems(c *gin.Context) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := orderItemCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing order items",
		}
	}

	var allOrderItems []bson.M

	if err := result.All(ctx, &allOrderItems); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return allOrderItems, nil
}

func ItemsByOrder(orderId string) (orderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{{Key: "order_id", Value: orderId}}},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food"},
			{Key: "localField", Value: "food_id"},
			{Key: "foreignField", Value: "food_id"},
			{Key: "as", Value: "food"},
		}},
	}

	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$food"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	lookupOrderStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "order"},
			{Key: "localField", Value: "order_id"},
			{Key: "foreignField", Value: "order_id"},
			{Key: "as", Value: "order"},
		}},
	}

	unwindOrderStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$order"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	lookupTableStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "table"},
			{Key: "localField", Value: "order.table_id"},
			{Key: "foreignField", Value: "table_id"},
			{Key: "as", Value: "table"},
		}},
	}

	unwindTableStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$table"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "food_name", Value: "$food.name"},
			{Key: "food_image", Value: "$food.food_image"},
			{Key: "table_number", Value: "$table.table_number"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "price", Value: "$food.price"},
			{Key: "quantity", Value: 1},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "order_id", Value: "$order_id"},
				{Key: "table_id", Value: "$table_id"},
				{Key: "table_number", Value: "$table_number"},
			}},
			{Key: "payment_due", Value: bson.D{
				{Key: "$sum", Value: "$amount"},
			}},
			{Key: "total_count", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
			{Key: "order_items", Value: bson.D{
				{Key: "$push", Value: "$$ROOT"},
			}},
		}},
	}

	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$lookup", Value: lookupStage}},
		{{Key: "$unwind", Value: unwindStage}},
		{{Key: "$lookup", Value: lookupOrderStage}},
		{{Key: "$unwind", Value: unwindOrderStage}},
		{{Key: "$lookup", Value: lookupTableStage}},
		{{Key: "$unwind", Value: unwindTableStage}},
		{{Key: "$project", Value: projectStage}},
		{{Key: "$group", Value: groupStage}},
		{{Key: "$project", Value: projectStage2}},
	})

	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := result.All(ctx, &orderItems); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return orderItems, nil

}

func GetOrderItem(c *gin.Context) (models.OrderItem, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("order_item_id")
	var orderItem models.OrderItem

	err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
		return models.OrderItem{}, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "error occured while listing ordered item",
		}
	}

	return orderItem, nil
}

func CreateOrderItem(c *gin.Context) (*mongo.InsertManyResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItemPack OrderItemPack
	var order models.Order

	if err := c.BindJSON(&orderItemPack); err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	orderItemsToBeInserted := []interface{}{}
	order.TableID = orderItemPack.TableID
	order_id, err := OrderItemOrderCreator(ctx, order)
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	for _, orderItem := range orderItemPack.OrderItems {
		orderItem.OrderID = order_id

		validationErr := validate.Struct(orderItem)

		if validationErr != nil {
			return nil, helpers.HttpError{
				Code:    http.StatusBadRequest,
				Message: validationErr.Error(),
			}
		}

		orderItem.ID = primitive.NewObjectID()
		orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.OrderItemID = orderItem.ID.Hex()
		var num = helpers.ToFixed(*orderItem.UnitPrice, 2)
		orderItem.UnitPrice = &num
		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
	if err != nil {
		return nil, helpers.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return insertedOrderItems, nil
}

func UpdateOrderItem(c *gin.Context) (*mongo.UpdateResult, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var orderItem models.OrderItem

	orderItemId := c.Param("order_item_id")

	filter := bson.M{"order_item_id": orderItemId}

	var updateObj primitive.D

	if orderItem.UnitPrice != nil {
		updateObj = append(updateObj, primitive.E{Key: "unit_price", Value: orderItem.UnitPrice})
	}

	if orderItem.Quantity != nil {
		updateObj = append(updateObj, primitive.E{Key: "quantity", Value: orderItem.Quantity})
	}

	if orderItem.FoodID != nil {
		updateObj = append(updateObj, primitive.E{Key: "food_id", Value: orderItem.FoodID})
	}

	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: orderItem.UpdatedAt})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(
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
			Message: "Order item update failed",
		}
	}

	return result, nil
}
