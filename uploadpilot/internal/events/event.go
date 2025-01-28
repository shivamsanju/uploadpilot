package events

import (
	"context"
	"sync"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

type UploadEventKey string

const (
	EventUploadStarted             UploadEventKey = "upload_started"
	EventUploadSkipped             UploadEventKey = "upload_skipped"
	EventUploadInProgress          UploadEventKey = "upload_in_progress"
	EventUploadFailed              UploadEventKey = "upload_failed"
	EventUploadComplete            UploadEventKey = "upload_complete"
	EventUploadCancelled           UploadEventKey = "upload_cancelled"
	EventUploadProcessing          UploadEventKey = "processing_upload"
	EventUploadProcessingFailed    UploadEventKey = "upload_processing_failed"
	EventUploadProcessingCancelled UploadEventKey = "upload_processing_canceled"
	EventUploadProcessed           UploadEventKey = "upload_processed"
	EventUploadDeleted             UploadEventKey = "upload_deleted"
)

type UploadEvent struct {
	Key     UploadEventKey
	Upload  models.Upload // this is a copy so no subscriber can change the upload
	Context context.Context
	Message string
	Error   error
}

type UploadEventBus struct {
	subscribers map[UploadEventKey][]chan UploadEvent
	lock        sync.RWMutex
}

var uploadEventOnce sync.Once
var uploadEventBus *UploadEventBus

func NewUploadEvent(ctx context.Context, Key UploadEventKey, upload *models.Upload, message string, err error) UploadEvent {
	return UploadEvent{
		Key:     Key,
		Upload:  *upload,
		Context: ctx,
		Message: message,
		Error:   err,
	}
}

func GetUploadEventBus() *UploadEventBus {
	uploadEventOnce.Do(func() {
		uploadEventBus = &UploadEventBus{
			subscribers: map[UploadEventKey][]chan UploadEvent{},
			lock:        sync.RWMutex{},
		}
	})

	return uploadEventBus
}

func (bus *UploadEventBus) Subscribe(key UploadEventKey, subscriber chan UploadEvent) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, exists := bus.subscribers[key]; !exists {
		bus.subscribers[key] = []chan UploadEvent{}
	}
	bus.subscribers[key] = append(bus.subscribers[key], subscriber)
}

func (bus *UploadEventBus) Publish(event UploadEvent) {
	infra.Log.Infof("publishing event %s", event.Key)
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	if subs, exists := bus.subscribers[event.Key]; exists {
		for _, sub := range subs {
			infra.Log.Warnf("publishing event %s to subscriber %p", event.Key, sub)
			go func(s chan UploadEvent) {
				s <- event
			}(sub)
		}
	} else {
		infra.Log.Warnf("no subscribers for event %s", event.Key)
	}
}

func (bus *UploadEventBus) Close() {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for _, subs := range bus.subscribers {
		for _, sub := range subs {
			close(sub)
		}
	}
	bus.subscribers = nil
}
