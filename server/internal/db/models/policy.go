package models

import "go.mongodb.org/mongo-driver/v2/bson"

type ImportPolicy struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string        `json:"name"`
	AllowedMimeTypes []string      `json:"allowedMimeTypes"`
	AllowedSources   []string      `json:"allowedSources"`
	MaxFileSize      int64         `json:"maxFileSize"`
	MaxFileCount     int64         `json:"maxFileCount"`
}
