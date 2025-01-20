package storage

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
)

func S3Client(ctx context.Context) (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		config.S3AccessKey,
		config.S3SecretKey,
		"",
	)

	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(config.S3Region),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awscfg), nil
}
