package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Helal struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name"`
	Age     int                `json:"age"`
	Surname string             `json:"surname"`
}

