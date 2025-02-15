package listeners

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	commonutils "github.com/uploadpilot/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/pubsub/pkg/events"
	"github.com/uploadpilot/uploadpilot/uploader/internal/infra"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/upload"
)

type StatusListener struct {
	uploadEvent    *events.UploadStatusEvent
	uploadLogEvent *events.UploadLogEvent
	uploadSvc      *upload.Service
}

func NewStatusListener(uploadSvc *upload.Service) *StatusListener {
	return &StatusListener{
		uploadEvent:    events.NewUploadStatusEvent(infra.RedisClient),
		uploadLogEvent: events.NewUploadLogEvent(infra.RedisClient),
		uploadSvc:      uploadSvc,
	}
}

func (l *StatusListener) Start() {
	defer commonutils.Recover()
	infra.Log.Info("starting upload status listener...")

	consumerGroup := "upload-status-listener"
	consumerKey := uuid.NewString()

	l.uploadEvent.Subscribe(consumerGroup, consumerKey, l.statusHandler)
}

func (l *StatusListener) statusHandler(msg *events.UploadEventMsg) error {
	ctx := context.Background()
	l.uploadLogEvent.Publish(
		msg.WorkspaceID,
		msg.UploadID,
		nil,
		nil,
		fmt.Sprintf("upload status changed to %s", msg.Status),
		string(models.UploadLogLevelInfo),
	)
	return l.uploadSvc.SetStatus(ctx, msg.UploadID, models.UploadStatus(msg.Status))
}
