package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/uploadpilot/uploadpilot/momentum/pkg/dsl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		Namespace: "test.o4ymi",
		HostPort:  "ap-southeast-1.aws.api.temporal.io:7233",
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(
					func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
						return invoker(
							metadata.AppendToOutgoingContext(ctx, "temporal-namespace", "test.o4ymi"),
							method,
							req,
							reply,
							cc,
							opts...,
						)
					},
				),
			},
		},
		Credentials: client.NewAPIKeyStaticCredentials("eyJhbGciOiJFUzI1NiIsImtpZCI6Ild2dHdhQSJ9.eyJhY2NvdW50X2lkIjoibzR5bWkiLCJhdWQiOlsidGVtcG9yYWwuaW8iXSwiZXhwIjoxNzQxODYyNTczLCJpc3MiOiJ0ZW1wb3JhbC5pbyIsImp0aSI6IkJqTGphSWRVV2tNTHZnajlxVHEyZFROYXhINmRUR254Iiwia2V5X2lkIjoiQmpMamFJZFVXa01MdmdqOXFUcTJkVE5heEg2ZFRHbngiLCJzdWIiOiI1OWU5MDY5OGE0YmE0ZGFlOGQ5NzdmZjYxMDUxMWUxMyJ9.NukpzRLqXRBRJs5Hc-uu3J-Z2SQ72rA3jFaHEEDLPRs8Ibu6AbF6krBs4i0Jemd0nXoTDHAPCumNOjBW57mCZw"),
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, "dsl", worker.Options{})

	w.RegisterWorkflow(dsl.SimpleDSLWorkflow)
	w.RegisterActivity(&dsl.SampleActivities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
