package listeners

import (
	"context"
	"sync"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadLogsListener struct {
	eventChan      chan *events.LogEvent
	uploadLogsRepo *db.UploadLogsRepo
	logBuffer      []interface{}
	flushInterval  time.Duration
	bufferSize     int
	mu             sync.Mutex
}

func NewUploadLogsListener(flushInterval time.Duration, bufferSize int) *UploadLogsListener {
	eventBus := events.GetLogEventBus()

	eventChan := make(chan *events.LogEvent)
	eventBus.Subscribe(eventChan)

	return &UploadLogsListener{
		eventChan:      eventChan,
		uploadLogsRepo: db.NewUploadLogsRepo(),
		logBuffer:      make([]interface{}, bufferSize),
		flushInterval:  flushInterval,
		bufferSize:     bufferSize,
		mu:             sync.Mutex{},
	}
}

func (l *UploadLogsListener) Start() {
	infra.Log.Info("starting upload logs listener...")
	go l.startBatchFlush()

	for event := range l.eventChan {
		infra.Log.Infof("processing log event %s", event.UploadID)

		uploadID, err := primitive.ObjectIDFromHex(event.UploadID)
		if err != nil {
			infra.Log.Errorf("invalid upload id in log: %s", err.Error())
			continue
		}

		workspaceID, err := primitive.ObjectIDFromHex(event.WorkspaceID)
		if err != nil {
			infra.Log.Errorf("invalid workspace id in log: %s", err.Error())
			continue
		}

		l.mu.Lock()

		l.logBuffer = append(l.logBuffer, &models.UploadLog{
			ID:          primitive.NewObjectID(),
			WorkspaceID: workspaceID,
			UploadID:    uploadID,
			Message:     event.Message,
			Level:       event.Level,
			Timestamp:   primitive.NewDateTimeFromTime(time.Now()),
		})

		if len(l.logBuffer) >= l.bufferSize {
			l.uploadLogsRepo.BatchAddLogs(context.Background(), l.logBuffer)
			l.logBuffer = nil
		}
		l.mu.Unlock()
	}
}

// startBatchFlush will flush logs every `flushInterval`
func (l *UploadLogsListener) startBatchFlush() {
	ticker := time.NewTicker(l.flushInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		if len(l.logBuffer) > 0 {
			infra.Log.Info("flushing logs...")
			l.uploadLogsRepo.BatchAddLogs(context.Background(), l.logBuffer)
			l.logBuffer = nil
		}
		l.mu.Unlock()
	}
}
