package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TemporalOpts struct {
	Namespace string
	HostPort  string
	APIKey    string
}

func NewTemporalClient(opts *TemporalOpts) (client.Client, error) {
	if opts.Namespace == "" || opts.HostPort == "" || opts.APIKey == "" {
		return nil, fmt.Errorf("temporal namespace, host port and api key are required")
	}

	c, err := client.Dial(client.Options{
		Logger:    log.DefaultLogger.Slog(),
		Namespace: opts.Namespace,
		HostPort:  opts.HostPort,
		ConnectionOptions: client.ConnectionOptions{
			GetSystemInfoTimeout: time.Second * 60,
			TLS:                  &tls.Config{},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(
					func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, grpcOpts ...grpc.CallOption) error {
						return invoker(
							metadata.AppendToOutgoingContext(ctx, "temporal-namespace", opts.Namespace),
							method,
							req,
							reply,
							cc,
							grpcOpts...,
						)
					},
				),
			},
		},
		Credentials: client.NewAPIKeyStaticCredentials(opts.APIKey),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create temporal client")
	}
	return c, nil
}
