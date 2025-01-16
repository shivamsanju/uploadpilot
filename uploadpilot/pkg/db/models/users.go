package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           primitive.ObjectID  `bson:"_id" json:"id" validate:"required"`
	FirstName    string              `bson:"firstName" json:"firstName" validate:"required,min=2,max=100"`
	LastName     string              `bson:"lastName" json:"lastName" validate:"required,min=2,max=100"`
	Password     string              `bson:"password" json:"-" validate:"required,min=6"`
	Email        string              `bson:"email" json:"email" validate:"email,required"`
	Phone        string              `bson:"phone" json:"phone" validate:"required"`
	Token        string              `bson:"token" json:"token"`
	RefreshToken string              `bson:"refreshToken" json:"refreshToken"`
	CreatedAt    primitive.Timestamp `bson:"createdAt" json:"createdAt"`
	UpdatedAt    primitive.Timestamp `bson:"updatedAt" json:"updatedAt"`
}
