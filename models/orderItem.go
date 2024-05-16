package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID 				primitive.ObjectID 	`bson:"_id"`
	Quantity		*int				`json:"quantity" validate:"required,gt=0"`	
	UnitPrice		*float64			`json:"unit_price" validate:"required"`
	CreatedAt		time.Time			`json:"created_at"`
	UpdatedAt		time.Time			`json:"updated_at"`
	FoodID			*string				`json:"food_id" validate:"required"`
	OrderItemID		string				`json:"order_item_id" validate:"required"`
	OrderID			string				`json:"order_id" validate:"required"`
}