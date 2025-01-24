package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UploadStatus string

const (
	UploadStatusStarted            UploadStatus = "Started"
	UploadStatusSkipped            UploadStatus = "Skipped"
	UploadStatusInProgress         UploadStatus = "In Progress"
	UploadStatusComplete           UploadStatus = "Uploaded"
	UploadStatusFailed             UploadStatus = "Failed"
	UploadStatusCancelled          UploadStatus = "Cancelled"
	UploadStatusProcessing         UploadStatus = "Processing"
	UploadStatusProcessingFailed   UploadStatus = "Processing Failed"
	UploadStatusProcessingComplete UploadStatus = "Processing Complete"
	UploadStatusDeleted            UploadStatus = "Deleted"
)

type Upload struct {
	ID             primitive.ObjectID     `bson:"_id" json:"id"`
	WorkspaceID    primitive.ObjectID     `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Status         UploadStatus           `bson:"status" json:"status" validate:"required"`
	Metadata       map[string]interface{} `bson:"metadata" json:"metadata"`
	StoredFileName string                 `bson:"storedFileName" json:"storedFileName" validate:"required"`
	Size           int64                  `bson:"size" json:"size" validate:"required"`
	URL            string                 `bson:"url" json:"url"`
	ProcesedURL    string                 `bson:"processedUrl" json:"processedUrl"`
	StartedAt      primitive.DateTime     `bson:"startedAt" json:"startedAt"`
	FinishedAt     primitive.DateTime     `bson:"finishedAt" json:"finishedAt"`
}

var UploadTerminalStates = []UploadStatus{
	UploadStatusSkipped,
	UploadStatusComplete,
	UploadStatusFailed,
	UploadStatusCancelled,
	UploadStatusProcessingFailed,
	UploadStatusProcessingComplete,
	UploadStatusDeleted,
}

var UploadNonTerminalStates = []UploadStatus{
	UploadStatusStarted,
	UploadStatusInProgress,
	UploadStatusProcessing,
}

var UploadAllStates = append(UploadTerminalStates, UploadNonTerminalStates...)
