package listeners

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
)

type StatusListener struct {
	uploadEb   *pubsub.EventBus[events.UploadEventMsg]
	logEb      *pubsub.EventBus[events.UploadLogEventMsg]
	uploadRepo *db.UploadRepo
}

func NewStatusListener() *StatusListener {
	return &StatusListener{
		uploadRepo: db.NewUploadRepo(),
		logEb:      events.NewUploadLogEventBus(config.EventBusRedisConfig, uuid.New().String()),
		uploadEb:   events.NewUploadStatusEvent(config.EventBusRedisConfig, uuid.New().String()),
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
	return l.uploadRepo.SetStatus(ctx, msg.UploadID, models.UploadStatus(msg.Status))
}
