package listeners

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type StorageListener struct {
	eventChan   chan events.UploadEvent
	uploadRepo  *db.UploadRepo
	logEventBus *events.LogEventBus
}

func NewStorageListener() *StorageListener {
	eventBus := events.GetUploadEventBus()

	eventChan := make(chan events.UploadEvent)
	eventBus.Subscribe(events.EventUploadComplete, eventChan)

	return &StorageListener{
		eventChan:   eventChan,
		uploadRepo:  db.NewUploadRepo(),
		logEventBus: events.GetLogEventBus(),
	}
}

func (l *StorageListener) Start() {
	defer utils.Recover()

	infra.Log.Info("starting upload storage listener...")
	for event := range l.eventChan {
		ctx := context.Background()
		uploadID := event.Upload.ID

		// Delete tus info file from S3 (no need to throw an error)
		go l.cleanUploadInfoFileRoutine(ctx, uploadID)

		url, objectName, err := l.generateUploadURL(ctx, uploadID)
		if err != nil {
			l.logEventBus.Publish(events.NewLogEvent(ctx, event.Upload.WorkspaceID, uploadID, "Failed to generate a link of the uploaded file", nil, nil, models.UploadLogLevelError))
			infra.Log.Errorf("failed to generate upload url: %s", err.Error())
			continue
		}

		patchMap := map[string]interface{}{
			"url":              url,
			"stored_file_name": objectName,
		}

		if err := l.uploadRepo.Patch(ctx, uploadID, patchMap); err != nil {
			l.logEventBus.Publish(events.NewLogEvent(ctx, event.Upload.WorkspaceID, uploadID, "Failed to save the link of the uploaded file", nil, nil, models.UploadLogLevelError))
			infra.Log.Errorf("failed to patch upload: %s", err.Error())
		}
	}
}

func (l *StorageListener) cleanUploadInfoFileRoutine(ctx context.Context, uploadID string) {
	defer utils.Recover()

	objectFileName := uploadID

	infoFileName := objectFileName + ".info"

	infra.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &config.S3BucketName, Key: &infoFileName})
}

func (l *StorageListener) generateUploadURL(ctx context.Context, uploadID string) (string, string, error) {
	objectFileName := uploadID

	infra.Log.Infof("objectFileName will be: %s", objectFileName)

	url, err := s3.NewPresignClient(infra.S3Client).PresignGetObject(
		ctx,
		&s3.GetObjectInput{Bucket: &config.S3BucketName, Key: &objectFileName},
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(7 * 24 * time.Hour)
		},
	)

	if err != nil {
		return "", "", err
	}

	return url.URL, objectFileName, nil
}
