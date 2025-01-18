package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID            primitive.ObjectID  `bson:"_id" json:"id" validate:"required"`
	UserID        string              `bson:"userId" json:"userId" validate:"required"`
	Email         string              `bson:"email" json:"email" validate:"email,required"`
	EmailVerified bool                `bson:"emailVerified" json:"emailVerified"`
	FirstName     string              `bson:"firstName" json:"firstName" validate:"required,min=2,max=100"`
	LastName      string              `bson:"lastName" json:"lastName" validate:"required,min=2,max=100"`
	Name          string              `bson:"name" json:"name" validate:"required"`
	NickName      string              `bson:"nickName" json:"nickName"`
	Phone         string              `bson:"phone" json:"phone" validate:"required"`
	Provider      string              `bson:"provider" json:"provider"`
	Description   string              `bson:"description" json:"description" validate:"max=500"`
	AvatarURL     string              `bson:"avatarUrl" json:"avatarUrl"`
	Location      string              `bson:"location" json:"location"`
	IsUserBanned  bool                `bson:"isUserBanned" json:"isUserBanned"`
	BanReason     string              `bson:"banReason" json:"banReason"`
	CreatedAt     primitive.Timestamp `bson:"createdAt" json:"createdAt"`
	UpdatedAt     primitive.Timestamp `bson:"updatedAt" json:"updatedAt"`
}
