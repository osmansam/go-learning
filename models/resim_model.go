package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Resim struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Adi string `json:"adi"`
	Resim string `json:"resim"`
}