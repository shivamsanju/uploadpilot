package listeners

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
)

type UploadLogsListener struct {
	logEb          *pubsub.EventBus[events.UploadLogEventMsg]
	uploadLogsRepo *db.UploadLogsRepo
	logBuffer      []*models.UploadLog
	flushInterval  time.Duration
	bufferSize     int
	mu             sync.Mutex
}

func NewUploadLogsListener(flushInterval time.Duration, bufferSize int) *UploadLogsListener {
	return &UploadLogsListener{
		logEb:          events.NewUploadLogEventBus(config.EventBusRedisConfig, uuid.New().String()),
		uploadLogsRepo: db.NewUploadLogsRepo(),
		logBuffer:      make([]*models.UploadLog, 0),
		flushInterval:  flushInterval,
		bufferSize:     bufferSize,
		mu:             sync.Mutex{},
	}
}

func (l *UploadLogsListener) Start() {
	defer commonutils.Recover()
	group := "upload-logs-listener"
	l.logEb.Subscribe(group, l.logHandler)
	go l.startBatchFlushInterval(l.flushInterval)
}

func (l *UploadLogsListener) logHandler(msg *events.UploadLogEventMsg) error {

	infra.Log.Infof("processing log event %s", msg.UploadID)
	l.mu.Lock()

	log := &models.UploadLog{
		WorkspaceID: msg.WorkspaceID,
		UploadID:    msg.UploadID,
		Message:     msg.Message,
		Level:       msg.Level,
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
		l.uploadLogsRepo.BatchAddLogs(context.Background(), l.logBuffer)
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
			l.uploadLogsRepo.BatchAddLogs(context.Background(), l.logBuffer)
			l.logBuffer = nil
		}
		l.mu.Unlock()
	}
}
