package listeners

import (
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/webhook"
)

type WebhookListener struct {
	eventChan  chan events.UploadEvent
	webhookSvc *webhook.WebhookService
}

func NewWebhookListener() *WebhookListener {
	eventBus := events.GetUploadEventBus()

	eventChan := make(chan events.UploadEvent)
	eventBus.Subscribe(events.EventUploadComplete, eventChan)

	return &WebhookListener{
		eventChan:  eventChan,
		webhookSvc: webhook.NewWebhookService(),
	}
}

func (l *WebhookListener) Start() {
	infra.Log.Info("starting webhook listener...")
	for event := range l.eventChan {
		_ = l.webhookSvc.TriggerWebhook(event.Upload, &event)
	}
}
