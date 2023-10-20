package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TimeStamps structure for common timestamp fields
type TimeStamps struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Token struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RefreshToken string             `json:"refresh_token"`
	Ip           string             `json:"ip"`
	UserAgent    string             `json:"user_agent"`
	IsValid      bool               `json:"is_valid"`
	User         primitive.ObjectID `json:"user"`
	TimeStamps   `json:",inline"`
}
