package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
)

type UploadEventMsg struct {
	WorkspaceID string
	UploadID    string
	Status      string
	Message     *string
	Error       error
}

func NewUploadEventMessage(workspaceID, uploadID, status string, message *string, err error) *UploadEventMsg {
	return &UploadEventMsg{WorkspaceID: workspaceID, UploadID: uploadID, Status: status, Message: message, Error: err}
}

func NewUploadStatusEvent(c *redis.Client, consumerKey string) *pubsub.EventBus[UploadEventMsg] {
	event := "up_upload_event"
	return pubsub.NewEventBus[UploadEventMsg](event, consumerKey, c)
}
