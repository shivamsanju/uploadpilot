package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id" validate:"required"`
	FirstName    string              `json:"firstName" validate:"required,min=2,max=100"`
	LastName     string              `json:"lastName" validate:"required,min=2,max=100"`
	Password     string              `json:"-" validate:"required,min=6"`
	Email        string              `json:"email" validate:"email,required"`
	Phone        string              `json:"phone" validate:"required"`
	Token        string              `json:"token"`
	RefreshToken string              `json:"refreshToken"`
	CreatedAt    primitive.Timestamp `json:"createdAt"`
	UpdatedAt    primitive.Timestamp `json:"updatedAt"`
}
