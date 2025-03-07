package clients

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsOpts struct {
	AccessKey string
	SecretKey string
	Region    string
}

func NewS3Client(opts *AwsOpts) (*s3.Client, error) {
	if opts.AccessKey == "" || opts.SecretKey == "" || opts.Region == "" {
		return nil, fmt.Errorf("s3 access key, secret key and region are required")
	}

	ctx := context.Background()
	creds := credentials.NewStaticCredentialsProvider(
		opts.AccessKey,
		opts.SecretKey,
		"",
	)

	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(opts.Region),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awscfg), nil
}

func NewLambdaClient(opts *AwsOpts) (*lambda.Client, error) {
	if opts.AccessKey == "" || opts.SecretKey == "" || opts.Region == "" {
		return nil, fmt.Errorf("s3 access key, secret key and region are required")
	}

	ctx := context.Background()
	creds := credentials.NewStaticCredentialsProvider(
		opts.AccessKey,
		opts.SecretKey,
		"",
	)

	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(opts.Region),
		awsconfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	return lambda.NewFromConfig(awscfg), nil
}
