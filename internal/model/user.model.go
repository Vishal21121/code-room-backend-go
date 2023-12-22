package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	IsLoggedIn   bool               `json:"isLoggedIn"`
	RefreshToken string             `json:"refreshToken"`
}
