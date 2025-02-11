package infra

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

var (
	Log            *zap.SugaredLogger
	S3Client       *s3.Client
	TemporalClient client.Client
)

func Init(s3Config *S3Config, temporalConfig *TemporalConfig) error {
	log, err := NewLogger()
	if err != nil {
		return err
	}
	Log = log

	if s3Config != nil {
		c, err := NewS3Client(s3Config)
		if err != nil {
			return err
		}
		S3Client = c
	}

	if temporalConfig != nil {
		tc, err := NewTemporalClient(temporalConfig)
		if err != nil {
			return err
		}
		TemporalClient = tc
	}

	return nil
}
