package listeners

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/upload"
)

type StatusListener struct {
	uploadEb  *pubsub.EventBus[events.UploadEventMsg]
	logEb     *pubsub.EventBus[events.UploadLogEventMsg]
	uploadSvc *upload.Service
}

func NewStatusListener(uploadSvc *upload.Service) *StatusListener {
	return &StatusListener{
		uploadSvc: uploadSvc,
		logEb:     events.NewUploadLogEventBus(infra.RedisClient, uuid.New().String()),
		uploadEb:  events.NewUploadStatusEvent(infra.RedisClient, uuid.New().String()),
	}
}

func (l *StatusListener) Start() {
	defer commonutils.Recover()
	infra.Log.Info("starting upload status listener...")
	group := "upload-status-listener"
	l.uploadEb.Subscribe(group, l.statusHandler)
}

func (l *StatusListener) statusHandler(msg *events.UploadEventMsg) error {
	ctx := context.Background()
	l.logEb.Publish(events.NewUploadLogEventMessage(
		msg.WorkspaceID,
		msg.UploadID,
		nil,
		nil,
		fmt.Sprintf("upload status changed to %s", msg.Status),
		models.UploadLogLevelInfo,
	))
	return l.uploadSvc.SetStatus(ctx, msg.UploadID, models.UploadStatus(msg.Status))
}
