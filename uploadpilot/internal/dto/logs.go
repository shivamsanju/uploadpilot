package dto

import (
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadLogNoIDs struct {
	ID          primitive.ObjectID    `bson:"_id" json:"id"`
	Level       models.UploadLogLevel `bson:"level" json:"level"`
	Timestamp   primitive.DateTime    `bson:"timestamp" json:"timestamp"`
	Message     string                `bson:"message" json:"message"`
	ProcessorID primitive.ObjectID    `bson:"processorId" json:"processorId"`
	TaskID      string                `bson:"taskId" json:"taskId"`
}
