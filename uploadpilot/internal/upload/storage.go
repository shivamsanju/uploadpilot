package upload

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/storage"
)

func (us *UploadService) GenerateS3URLFromUploadID(ctx context.Context, uploadID string) (string, string, error) {
	client, err := storage.S3Client(ctx)
	if err != nil {
		return "", "", err
	}

	objectFileName := uploadID
	if len(objectFileName) > 32 {
		objectFileName = objectFileName[:32]
	}

	infoFileName := objectFileName + ".info"

	// Delete tus info file from S3 (no need to throw an error)
	client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &config.S3BucketName, Key: &infoFileName})

	url, err := s3.NewPresignClient(client).PresignGetObject(
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
