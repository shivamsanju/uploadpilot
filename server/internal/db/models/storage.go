package models

import "go.mongodb.org/mongo-driver/v2/bson"

type StorageUnit struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ConnectorID bson.ObjectID `json:"connectorId" bson:"connectorId"`
	Bucket      string        `json:"bucket"`
}
