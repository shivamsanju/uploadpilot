package clients

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/core/pkg/vault"
	"go.temporal.io/sdk/client"
)

type Clients struct {
	RedisClient    *redis.Client
	S3Client       *s3.Client
	LambdaClient   *lambda.Client
	TemporalClient client.Client
	KMSClient      vault.KMS
}

type ClientOpts struct {
	RedisOpts    *RedisOpts
	S3Opts       *AwsOpts
	LambdaOpts   *AwsOpts
	TemporalOpts *TemporalOpts
	KMSOpts      *KMSOpts
}

func NewAppClients(opts *ClientOpts) (*Clients, error) {
	c := &Clients{}

	if opts.RedisOpts != nil {
		redisClient, err := NewRedisClient(opts.RedisOpts)
		if err != nil {
			return nil, err
		}
		c.RedisClient = redisClient
	} else {
		log.Warn().Msg("redis client not initialized")
	}

	if opts.S3Opts != nil {
		s3Client, err := NewS3Client(opts.S3Opts)
		if err != nil {
			return nil, err
		}
		c.S3Client = s3Client
	} else {
		log.Warn().Msg("S3 client not initialized")
	}

	if opts.LambdaOpts != nil {
		lambdaClient, err := NewLambdaClient(opts.LambdaOpts)
		if err != nil {
			return nil, err
		}
		c.LambdaClient = lambdaClient
	} else {
		log.Warn().Msg("lambda client not initialized")
	}

	if opts.TemporalOpts != nil {
		temporalClient, err := NewTemporalClient(opts.TemporalOpts)
		if err != nil {
			return nil, err
		}
		c.TemporalClient = temporalClient
	} else {
		log.Warn().Msg("temporal client not initialized")
	}

	if opts.KMSOpts != nil {
		kms, err := NewKMSClient(opts.KMSOpts)
		if err != nil {
			return nil, err
		}
		c.KMSClient = kms
	} else {
		log.Warn().Msg("KMS client not initialized")
	}

	return c, nil
}
