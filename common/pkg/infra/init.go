package infra

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

var (
	Log            *zap.SugaredLogger
	S3Client       *s3.Client
	TemporalClient client.Client
	RedisClient    *redis.Client
)

type InfraOpts struct {
	S3Opts       *S3Options
	TemporalOpts *TemporalOptions
	RedisOpts    *redis.Options
}

func Init(opts *InfraOpts) error {
	log, err := NewLogger()
	if err != nil {
		return err
	}
	Log = log

	if opts.S3Opts != nil {
		c, err := NewS3Client(opts.S3Opts)
		if err != nil {
			return err
		}
		S3Client = c
	} else {
		log.Warn("S3 client not initialized")
	}

	if opts.TemporalOpts != nil {
		tc, err := NewTemporalClient(opts.TemporalOpts)
		if err != nil {
			return err
		}
		TemporalClient = tc
	} else {
		log.Warn("Temporal client not initialized")
	}

	if opts.RedisOpts != nil {
		rc := redis.NewClient(opts.RedisOpts)
		cmd := rc.Ping(context.Background())
		if err := cmd.Err(); err != nil {
			return err
		}
		RedisClient = rc
	} else {
		log.Warn("Redis client not initialized")
	}

	return nil
}
