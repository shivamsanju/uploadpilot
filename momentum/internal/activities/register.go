package activities

import (
	"context"

	activitycatalog "github.com/uploadpilot/go-core/common/activitycatalog"
	"github.com/uploadpilot/momentum/internal/activities/img"
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
	ar.register(activitycatalog.ExtractPDFContentV1_0.Name, pdfActivities.ExtractContentFromPDF)

	// Register webhook activities
	webhookActivities := &WebhookActivities{}
	ar.register(activitycatalog.WebhookV_1_0.Name, webhookActivities.SendWebhook)

	// Register image activities
	ar.register(activitycatalog.ImageResizeV1_0.Name, img.ResizeImage)
	ar.register(activitycatalog.ImageConvertToPngV1_0.Name, img.ConvertImageToPng)
	ar.register(activitycatalog.ImageConvertToBmpV1_0.Name, img.ConvertImageToBmp)
	ar.register(activitycatalog.ImageConvertToJpegV1_0.Name, img.ConvertImageToJpeg)
	ar.register(activitycatalog.ImageAddWatermarkActivityV1_0.Name, img.ApplyTextWatermark)
	ar.register(activitycatalog.ImageMetadataExtractionActivityV1_0.Name, img.ExtractMetadataFromImage)
	ar.register(activitycatalog.ImageBlurActivityV1_0.Name, img.BlurImage)
}
