package listeners

import (
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type StatusListener struct {
	eventChan  chan events.UploadEvent
	uploadRepo *db.UploadRepo
}

var EventStatusMap = map[events.UploadEventKey]models.UploadStatus{
	events.EventUploadStarted:             models.UploadStatusStarted,
	events.EventUploadInProgress:          models.UploadStatusInProgress,
	events.EventUploadSkipped:             models.UploadStatusSkipped,
	events.EventUploadComplete:            models.UploadStatusComplete,
	events.EventUploadFailed:              models.UploadStatusFailed,
	events.EventUploadCancelled:           models.UploadStatusCancelled,
	events.EventUploadDeleted:             models.UploadStatusDeleted,
	events.EventUploadProcessing:          models.UploadStatusProcessing,
	events.EventUploadProcessingFailed:    models.UploadStatusProcessingFailed,
	events.EventUploadProcessed:           models.UploadStatusProcessingComplete,
	events.EventUploadProcessingCancelled: models.UploadStatusProcessingCancelled,
}

func NewStatusListener() *StatusListener {
	eventBus := events.GetUploadEventBus()

	eventChan := make(chan events.UploadEvent)
	for key := range EventStatusMap {
		eventBus.Subscribe(key, eventChan)
	}

	return &StatusListener{
		eventChan:  eventChan,
		uploadRepo: db.NewUploadRepo(),
	}
}

func (l *StatusListener) Start() {
	defer utils.Recover()
	infra.Log.Info("starting upload status listener...")
	for event := range l.eventChan {
		infra.Log.Infof("processing upload event %s", event.Key)
		status, ok := EventStatusMap[event.Key]
		if !ok {
			infra.Log.Warn("skipping unknown event key: %s", event.Key)
			continue
		}
		_ = l.uploadRepo.SetStatus(event.Context, event.Upload.ID.Hex(), status)
	}
}
