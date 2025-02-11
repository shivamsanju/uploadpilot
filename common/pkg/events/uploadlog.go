package events

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
)

type UploadLogEventMsg struct {
	Level       models.UploadLogLevel
	Message     string
	WorkspaceID string
	UploadID    string
	ProcessorID *string
	TaskID      *string
}

func NewUploadLogEventMessage(workspaceID, uploadID string, processorID, taskID *string, message string, level models.UploadLogLevel) *UploadLogEventMsg {
	return &UploadLogEventMsg{
		WorkspaceID: workspaceID,
		UploadID:    uploadID,
		ProcessorID: processorID,
		TaskID:      taskID,
		Message:     message,
		Level:       level,
	}
}

func NewUploadLogEventBus(c *pubsub.RedisConfig, consumerKey string) *pubsub.EventBus[UploadLogEventMsg] {
	event := "ul_upload_log"
	return pubsub.NewEventBus[UploadLogEventMsg](event, consumerKey, c)
}
