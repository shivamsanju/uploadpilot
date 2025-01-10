package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataStore struct {
	ID          primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Name        string             `json:"name" validate:"required,min=2,max=100"`
	ConnectorID primitive.ObjectID `json:"connectorId" validate:"required"`
	Bucket      string             `json:"bucket" validate:"required"`
	CreatedAt   primitive.DateTime `json:"createdAt"`
	CreatedBy   string             `json:"createdBy"`
	UpdatedAt   primitive.DateTime `json:"updatedAt"`
	UpdatedBy   string             `json:"updatedBy"`
}

type DataStoreWithConnector struct {
	*DataStore
	Connector *StorageConnector `json:"connector"`
}
