package events

import (
	"context"
	"sync"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

type LogEvent struct {
	Context     context.Context
	WorkspaceID string
	UploadID    string
	Message     string
	Level       models.UploadLogLevel
}

type LogEventBus struct {
	subscribers []chan LogEvent
	lock        sync.RWMutex
}

var logEventBus *LogEventBus
var logEventOnce sync.Once

func NewLogEvent(ctx context.Context, workspaceID, uploadID, message string, level models.UploadLogLevel) LogEvent {
	return LogEvent{Context: ctx, WorkspaceID: workspaceID, UploadID: uploadID, Message: message, Level: level}
}

func GetLogEventBus() *LogEventBus {
	logEventOnce.Do(func() {
		logEventBus = &LogEventBus{
			subscribers: make([]chan LogEvent, 0),
			lock:        sync.RWMutex{},
		}
	})
	return logEventBus
}

func (bus *LogEventBus) Subscribe(sub chan LogEvent) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	bus.subscribers = append(bus.subscribers, sub)
}

func (bus *LogEventBus) Publish(event LogEvent) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	for _, sub := range bus.subscribers {
		infra.Log.Infof("publishing log event to subscriber %p", sub)
		go func(s chan LogEvent) {
			s <- event
		}(sub)
	}
}

func (bus *LogEventBus) Close() {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for _, sub := range bus.subscribers {
		close(sub)
	}
	bus.subscribers = nil
}
