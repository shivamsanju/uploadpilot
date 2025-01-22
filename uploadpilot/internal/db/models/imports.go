package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	TimeStamp primitive.DateTime `bson:"timestamp" json:"timestamp"`
	Message   string             `bson:"message" json:"message"`
}

type ImportStatus string

const (
	ImportStatusInProgress ImportStatus = "In Progress"
	ImportStatusFailed     ImportStatus = "Failed"
	ImportStatusSuccess    ImportStatus = "Success"
	ImportStatusCancelled  ImportStatus = "Cancelled"
	ImportStatusDeleted    ImportStatus = "Deleted"
)

type Import struct {
	ID             primitive.ObjectID     `bson:"_id" json:"id"`
	WorkspaceID    primitive.ObjectID     `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Status         ImportStatus           `bson:"status" json:"status" validate:"required"`
	UploadID       string                 `bson:"uploadId" json:"uploadId" validate:"required"`
	Metadata       map[string]interface{} `bson:"metadata" json:"metadata"`
	StoredFileName string                 `bson:"storedFileName" json:"storedFileName" validate:"required"`
	URL            string                 `bson:"url" json:"url"`
	Size           int64                  `bson:"size" json:"size" validate:"required"`
	StartedAt      primitive.DateTime     `bson:"startedAt" json:"startedAt"`
	FinishedAt     primitive.DateTime     `bson:"finishedAt" json:"finishedAt"`
	Logs           []Log                  `bson:"logs" json:"logs"`
}
