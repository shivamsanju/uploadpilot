package infra

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

var (
	Log      *zap.SugaredLogger
	S3Client *s3.Client
)

type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
}

func Init(s3Config *S3Config) error {
	log, err := NewLogger()
	if err != nil {
		return err
	}
	Log = log

	c, err := NewS3Client(s3Config.AccessKey, s3Config.SecretKey, s3Config.Region)
	if err != nil {
		return err
	}
	S3Client = c

	return nil
}
