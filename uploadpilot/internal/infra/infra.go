package infra

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"go.uber.org/zap"
)

var (
	Log      *zap.SugaredLogger
	Validate = validator.New()
	S3Client *s3.Client
)

func Init() error {
	log, err := NewLogger()
	if err != nil {
		return err
	}
	Log = log

	c, err := NewS3Client()
	if err != nil {
		return err
	}
	S3Client = c

	return nil
}
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar, nil
}

func NewS3Client() (*s3.Client, error) {
	ctx := context.Background()
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
	client := s3.NewFromConfig(awscfg)
	return client, nil
}
