package listeners

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	commonutils "github.com/uploadpilot/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/pubsub/pkg/events"
	"github.com/uploadpilot/uploadpilot/uploader/internal/infra"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/upload"
)

type UploadLogsListener struct {
	uploadLogEvent *events.UploadLogEvent
	uploadSvc      *upload.Service
	logBuffer      []*models.UploadLog
	flushInterval  time.Duration
	bufferSize     int
	mu             sync.Mutex
}

func NewUploadLogsListener(flushInterval time.Duration, bufferSize int, uploadSvc *upload.Service) *UploadLogsListener {
	return &UploadLogsListener{
		uploadSvc:      uploadSvc,
		uploadLogEvent: events.NewUploadLogEvent(infra.RedisClient),
		logBuffer:      make([]*models.UploadLog, 0),
		flushInterval:  flushInterval,
		bufferSize:     bufferSize,
		mu:             sync.Mutex{},
	}
}

func (l *UploadLogsListener) Start() {
	defer commonutils.Recover()

	consumerGroup := "upload-logs-listener"
	consumerKey := uuid.NewString()
	l.uploadLogEvent.Subscribe(consumerGroup, consumerKey, l.logHandler)

	go l.startBatchFlushInterval(l.flushInterval)
}

func (l *UploadLogsListener) logHandler(msg *events.UploadLogEventMsg) error {

	infra.Log.Infof("processing log event %s", msg.UploadID)
	l.mu.Lock()

	log := &models.UploadLog{
		WorkspaceID: msg.WorkspaceID,
		UploadID:    msg.UploadID,
		Message:     msg.Message,
		Level:       models.UploadLogLevel(msg.Level),
		Timestamp:   time.Now(),
	}

	if msg.ProcessorID != nil {
		log.ProcessorID = msg.ProcessorID
	}

	if msg.TaskID != nil {
		log.TaskID = msg.TaskID
	}

	l.logBuffer = append(l.logBuffer, log)

	if len(l.logBuffer) >= l.bufferSize {
		l.uploadSvc.BatchAddLogs(context.Background(), l.logBuffer)
		l.logBuffer = nil
	}

	l.mu.Unlock()
	return nil
}

// startBatchFlush will flush logs every `flushInterval`
func (l *UploadLogsListener) startBatchFlushInterval(flushInterval time.Duration) {
	defer commonutils.Recover()

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		if len(l.logBuffer) > 0 {
			infra.Log.Info("interval flushing logs...")
			l.uploadSvc.BatchAddLogs(context.Background(), l.logBuffer)
			l.logBuffer = nil
		}
		l.mu.Unlock()
	}
}
