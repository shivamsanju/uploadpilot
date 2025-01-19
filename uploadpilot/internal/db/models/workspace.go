package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkspaceUser struct {
	UserID string   `bson:"userId" json:"userId" validate:"required"`
	Role   UserRole `bson:"role" json:"role" validate:"required"`
}

type Workspace struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Users     []WorkspaceUser    `bson:"users" json:"users"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
}
