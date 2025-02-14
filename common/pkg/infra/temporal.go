package infra

import (
	"context"
	"crypto/tls"
	"log"

	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TemporalOptions struct {
	Namespace string
	HostPort  string
	APIKey    string
}

func NewTemporalClient(opts *TemporalOptions) (client.Client, error) {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		Namespace: opts.Namespace,
		HostPort:  opts.HostPort,
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{},
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
		log.Fatalln("Unable to create Temporal client.", err)
	}
	return c, nil
}
