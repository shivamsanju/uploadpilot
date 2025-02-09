package listeners

import (
	"context"
	"sync"
	"time"

	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type UploadLogsListener struct {
	eventChan      chan *events.LogEvent
	uploadLogsRepo *db.UploadLogsRepo
	logBuffer      []*models.UploadLog
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
		logBuffer:      make([]*models.UploadLog, 0),
		flushInterval:  flushInterval,
		bufferSize:     bufferSize,
		mu:             sync.Mutex{},
	}
}

func (l *UploadLogsListener) Start() {
	defer utils.Recover()

	infra.Log.Info("starting upload logs listener...")
	go l.startBatchFlush()

	for event := range l.eventChan {
		infra.Log.Infof("processing log event %s", event.UploadID)

		l.mu.Lock()

		log := &models.UploadLog{
			WorkspaceID: event.WorkspaceID,
			UploadID:    event.UploadID,
			Message:     event.Message,
			Level:       event.Level,
			Timestamp:   time.Now(),
		}

		if event.ProcessorID != nil {
			log.ProcessorID = event.ProcessorID
		}

		if event.TaskID != nil {
			log.TaskID = event.TaskID
		}

		l.logBuffer = append(l.logBuffer, log)

		if len(l.logBuffer) >= l.bufferSize {
			l.uploadLogsRepo.BatchAddLogs(context.Background(), l.logBuffer)
			l.logBuffer = nil
		}

		l.mu.Unlock()
	}
}

// startBatchFlush will flush logs every `flushInterval`
func (l *UploadLogsListener) startBatchFlush() {
	defer utils.Recover()

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
