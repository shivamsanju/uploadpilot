package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/go-core/pubsub/pkg/bus"
)

type UploadLogEventMsg struct {
	Level       string
	Message     string
	WorkspaceID string
	UploadID    string
	ProcessorID *string
	TaskID      *string
}

type UploadLogEvent struct {
	bus *bus.EventBus[UploadLogEventMsg]
}

func NewUploadLogEvent(c *redis.Client) *UploadLogEvent {
	event := "ul_upload_log"
	bus := bus.NewEventBus[UploadLogEventMsg](event, c)
	return &UploadLogEvent{bus: bus}
}

func (e *UploadLogEvent) Subscribe(consumerGroup string, consumerKey string, handler func(msg *UploadLogEventMsg) error) error {
	return e.bus.Subscribe(consumerGroup, consumerKey, handler)
}

func (e *UploadLogEvent) Publish(workspaceID, uploadID string, processorID, taskID *string, message string, level string) error {
	return e.bus.Publish(&UploadLogEventMsg{
		WorkspaceID: workspaceID,
		UploadID:    uploadID,
		ProcessorID: processorID,
		TaskID:      taskID,
		Message:     message,
		Level:       level,
	})
}
