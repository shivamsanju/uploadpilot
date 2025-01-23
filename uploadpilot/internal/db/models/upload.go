package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	TimeStamp primitive.DateTime `bson:"timestamp" json:"timestamp"`
	Message   string             `bson:"message" json:"message"`
}

type UploadStatus string

const (
	UploadStatusInProgress UploadStatus = "In Progress"
	UploadStatusFailed     UploadStatus = "Failed"
	UploadStatusSuccess    UploadStatus = "Success"
	UploadStatusCancelled  UploadStatus = "Cancelled"
	UploadStatusDeleted    UploadStatus = "Deleted"
)

type Upload struct {
	ID             primitive.ObjectID     `bson:"_id" json:"id"`
	WorkspaceID    primitive.ObjectID     `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Status         UploadStatus           `bson:"status" json:"status" validate:"required"`
	Metadata       map[string]interface{} `bson:"metadata" json:"metadata"`
	StoredFileName string                 `bson:"storedFileName" json:"storedFileName" validate:"required"`
	URL            string                 `bson:"url" json:"url"`
	Size           int64                  `bson:"size" json:"size" validate:"required"`
	StartedAt      primitive.DateTime     `bson:"startedAt" json:"startedAt"`
	FinishedAt     primitive.DateTime     `bson:"finishedAt" json:"finishedAt"`
	Logs           []Log                  `bson:"logs" json:"logs"`
}
