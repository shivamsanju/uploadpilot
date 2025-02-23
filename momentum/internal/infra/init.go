package infra

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/phuslu/log"
	"go.temporal.io/sdk/client"
)

var (
	S3Client       *s3.Client
	TemporalClient client.Client
)

type InfraOpts struct {
	S3Opts       *S3Options
	TemporalOpts *TemporalOptions
}

func Init(opts *InfraOpts) error {
	if opts.S3Opts != nil {
		c, err := NewS3Client(opts.S3Opts)
		if err != nil {
			return err
		}
		S3Client = c
	} else {
		log.Warn().Msg("S3 client not initialized")
	}

	if opts.TemporalOpts != nil {
		tc, err := NewTemporalClient(opts.TemporalOpts)
		if err != nil {
			return err
		}
		TemporalClient = tc
	} else {
		log.Warn().Msg("Temporal client not initialized")
	}

	return nil
}
