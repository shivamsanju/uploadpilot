package worker

import (
	"context"

	"github.com/uploadpilot/core/internal/workflow/activities"
	"github.com/uploadpilot/core/internal/workflow/catalog"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type ActivityRegistry struct {
	worker worker.Worker
}

func NewActivityRegistry(worker worker.Worker) *ActivityRegistry {
	return &ActivityRegistry{
		worker: worker,
	}
}

func (ar *ActivityRegistry) register(
	activityName string,
	activityFunc func(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error),
) {
	ar.worker.RegisterActivityWithOptions(
		activityFunc,
		activity.RegisterOptions{
			Name: activityName,
		},
	)
}

func (ar *ActivityRegistry) RegisterActivities(worker worker.Worker) {
	// Register all activities

	// Register pdf activities
	pdfActivities := &activities.PdfActivities{}
	ar.register(catalog.ExtractPDFContentV1_0.Name, pdfActivities.ExtractContentFromPDF)

	// Register webhook activities
	webhookActivities := &activities.WebhookActivities{}
	ar.register(catalog.WebhookV_1_0.Name, webhookActivities.SendWebhook)

}
