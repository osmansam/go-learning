package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct{
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Password string `json:"password"`
	Email string `json:"email"`
	Role string `json:"role"`
	Username string `json:"username"`
}