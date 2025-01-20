package hooks

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mitchellh/mapstructure"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateImportMetadata(hook tusd.HookEvent, s3Client *s3.Client, imp *models.Import) error {
	imp.ID = primitive.NewObjectID()
	imp.Size = hook.Upload.Size
	imp.StartedAt = primitive.NewDateTimeFromTime(time.Now())
	imp.Status = models.ImportStatusInProgress
	imp.Logs = []models.Log{{
		Message:   "Import started",
		TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
	}}

	// Extract Metadata
	var metadata map[string]interface{}
	err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
	if err != nil {
		return err
	}
	imp.Metadata = metadata

	// Add workspaceId
	headers := hook.HTTPRequest.Header
	wsID := headers.Get("workspaceId")
	if len(wsID) == 0 {
		return fmt.Errorf("missing workspaceId in header")
	}
	workspaceID, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		return fmt.Errorf("invalid workspaceId: %w", err)
	}
	imp.WorkspaceID = workspaceID

	// Add Object Name
	objectFileName := hook.Upload.ID
	if len(objectFileName) > 32 {
		objectFileName = objectFileName[:32]
	}
	imp.StoredFileName = objectFileName

	// Add Presigned URL
	URL, err := getPresignedURL(hook.Context, s3Client, objectFileName)
	if err != nil {
		return err
	}
	imp.URL = URL
	imp.Size = hook.Upload.Size
	return nil
}

func getPresignedURL(ctx context.Context, client *s3.Client, objectKey string) (string, error) {
	url, err := s3.NewPresignClient(client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &config.S3BucketName,
		Key:    &objectKey,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(7 * 24 * time.Hour)
	})

	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	infra.Log.Infof("presigned url -> %s", url.URL)
	return url.URL, nil
}

func RemoveInfoFile(hook tusd.HookEvent, s3Client *s3.Client, filename string) {
	objectFileName := hook.Upload.ID
	if len(objectFileName) > 32 {
		objectFileName = objectFileName[:32]
	}
	infoFileName := objectFileName + ".info"

	_, err := s3Client.DeleteObject(hook.Context, &s3.DeleteObjectInput{Bucket: &config.S3BucketName, Key: &infoFileName})
	if err != nil {
		infra.Log.Infof("failed to remove tus upload info file -> %+v", objectFileName)
	}
}
