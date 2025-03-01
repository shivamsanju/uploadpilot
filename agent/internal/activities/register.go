package activities

import (
	"context"

	"github.com/uploadpilot/agent/internal/activities/img"
	"github.com/uploadpilot/go-common/workflow/catalog"
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
	pdfActivities := &PdfActivities{}
	ar.register(catalog.ExtractPDFContentV1_0.Name, pdfActivities.ExtractContentFromPDF)

	// Register webhook activities
	webhookActivities := &WebhookActivities{}
	ar.register(catalog.WebhookV_1_0.Name, webhookActivities.SendWebhook)

	// Register image activities
	ar.register(catalog.ImageResizeV1_0.Name, img.ResizeImage)
	ar.register(catalog.ImageConvertToPngV1_0.Name, img.ConvertImageToPng)
	ar.register(catalog.ImageConvertToBmpV1_0.Name, img.ConvertImageToBmp)
	ar.register(catalog.ImageConvertToJpegV1_0.Name, img.ConvertImageToJpeg)
	ar.register(catalog.ImageAddWatermarkActivityV1_0.Name, img.ApplyTextWatermark)
	ar.register(catalog.ImageMetadataExtractionActivityV1_0.Name, img.ExtractMetadataFromImage)
	ar.register(catalog.ImageBlurActivityV1_0.Name, img.BlurImage)
}
