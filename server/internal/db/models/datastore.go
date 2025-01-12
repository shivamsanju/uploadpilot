package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataStore struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	ConnectorID primitive.ObjectID `bson:"connectorId" json:"connectorId" validate:"required"`
	Bucket      string             `bson:"bucket" json:"bucket" validate:"required"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
}
