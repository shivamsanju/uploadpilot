package activities

import (
	"github.com/uploadpilot/go-core/common/tasks"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

func RegisterActivities(worker worker.Worker) {
	// Register all activities
	RegisterPdfActivities(worker)
	RegisterWebhookActivities(worker)
}

func RegisterPdfActivities(worker worker.Worker) {
	// Register all activities
	pdfActivities := &PdfActivities{}

	worker.RegisterActivityWithOptions(
		pdfActivities.ExtractContentFromPDF,
		activity.RegisterOptions{
			Name: tasks.ExtractPDFContentTask.Name,
		},
	)

}

func RegisterWebhookActivities(worker worker.Worker) {
	// Register all activities
	webhookActivities := &WebhookActivities{}
	worker.RegisterActivityWithOptions(
		webhookActivities.SendWebhook,
		activity.RegisterOptions{
			Name: tasks.WebhookTask.Name,
		})
}
