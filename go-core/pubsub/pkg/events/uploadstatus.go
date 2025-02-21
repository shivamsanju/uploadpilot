package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/go-core/pubsub/pkg/bus"
)

type UploadEventMsg struct {
	WorkspaceID string
	UploadID    string
	Status      string
	Message     *string
	Error       error
}

type UploadStatusEvent struct {
	bus *bus.EventBus[UploadEventMsg]
}

func NewUploadStatusEvent(c *redis.Client) *UploadStatusEvent {
	event := "up_upload_event"
	b := bus.NewEventBus[UploadEventMsg](event, c)

	return &UploadStatusEvent{
		bus: b,
	}
}

func (e *UploadStatusEvent) Subscribe(consumerGroup string, consumerKey string, handler func(msg *UploadEventMsg) error) error {
	return e.bus.Subscribe(consumerGroup, consumerKey, handler)
}

func (e *UploadStatusEvent) Publish(workspaceID, uploadID, status string, message *string, err error) error {
	return e.bus.Publish(&UploadEventMsg{
		WorkspaceID: workspaceID,
		UploadID:    uploadID,
		Status:      status,
		Message:     message,
		Error:       err,
	})
}

func NewUploadEventMessage(workspaceID, uploadID, status string, message *string, err error) *UploadEventMsg {
	return &UploadEventMsg{WorkspaceID: workspaceID, UploadID: uploadID, Status: status, Message: message, Error: err}
}
