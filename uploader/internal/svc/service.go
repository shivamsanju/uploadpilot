package svc

import (
	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
)

type Service struct {
	configRepo    *db.WorkspaceConfigRepo
	uploadRepo    *db.UploadRepo
	procRepo      *db.ProcessorRepo
	workspaceRepo *db.WorkspaceRepo
	logEventBus   *pubsub.EventBus[events.UploadLogEventMsg]
	uploadEb      *pubsub.EventBus[events.UploadEventMsg]
}

func NewService() *Service {
	return &Service{
		configRepo:    db.NewWorkspaceConfigRepo(),
		uploadRepo:    db.NewUploadRepo(),
		workspaceRepo: db.NewWorkspaceRepo(),
		uploadEb:      events.NewUploadStatusEvent(config.EventBusRedisConfig, uuid.New().String()),
		logEventBus:   events.NewUploadLogEventBus(config.EventBusRedisConfig, uuid.New().String()),
	}
}
