package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadLogLevel string

const (
	UploadLogLevelInfo  UploadLogLevel = "info"
	UploadLogLevelWarn  UploadLogLevel = "warn"
	UploadLogLevelError UploadLogLevel = "error"
)

type UploadLog struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	WorkspaceID primitive.ObjectID `bson:"workspaceId" json:"workspaceId"`
	UploadID    primitive.ObjectID `bson:"uploadId" json:"uploadId"`
	Level       UploadLogLevel     `bson:"level" json:"level"`
	Timestamp   primitive.DateTime `bson:"timestamp" json:"timestamp"`
	Message     string             `bson:"message" json:"message"`
}
