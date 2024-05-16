package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID 			primitive.ObjectID 	`bson:"_id"`
	Text		string				`json:"text"`
	Title		string				`json:"title"`
	CreatedAt	time.Time			`json:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at"`
	NoteID		string				`json:"note_id"`
}