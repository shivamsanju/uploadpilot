package catalog

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
	"github.com/uploadpilot/uploadpilot/internal/storage"
)

func (c *hooksCatalogService) AddUploadURLHookFunc(ctx context.Context, input *hooks.HookInput, continueOnError bool) *hooks.HookResponse {
	client, err := storage.S3Client(ctx)
	if err != nil {
		return c.updateImportAndReturnErrorResponse(ctx, input, err, continueOnError)
	}
	objectFileName := input.TusdHook.Upload.ID
	if len(objectFileName) > 32 {
		objectFileName = objectFileName[:32]
	}

	input.Import.StoredFileName = objectFileName
	// Delete tus info file from S3 (no need to throw an error)
	client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &config.S3BucketName, Key: &objectFileName})

	url, err := s3.NewPresignClient(client).PresignGetObject(
		ctx,
		&s3.GetObjectInput{Bucket: &config.S3BucketName, Key: &objectFileName},
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(7 * 24 * time.Hour)
		},
	)
	if err != nil {
		return c.updateImportAndReturnErrorResponse(ctx, input, err, continueOnError)
	}

	input.Import.URL = url.URL
	return c.updateImportAndReturnSuccessResponse(ctx, input, "URL added to the import successfully")
}
