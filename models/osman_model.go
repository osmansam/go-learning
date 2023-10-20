package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Osman struct{
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name"`
	Helal primitive.ObjectID `json:"helal"`
	HelalDetails Helal `json:"helal_details,omitempty" bson:"helal_details,omitempty"`

}