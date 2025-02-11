package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"

	"github.com/pborman/uuid"
	"github.com/uploadpilot/uploadpilot/momentum/pkg/dsl"
	"github.com/uploadpilot/uploadpilot/momentum/pkg/validation"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v3"
)

func main() {
	var dslConfig string
	flag.StringVar(&dslConfig, "dslConfig", "pkg/dsl/workflow1.yaml", "dslConfig specify the yaml file for the dsl workflow.")
	flag.Parse()

	data, err := os.ReadFile(dslConfig)
	if err != nil {
		log.Fatalln("failed to load dsl config file", err)
	}
	var dslWorkflow dsl.Workflow
	if err := yaml.Unmarshal(data, &dslWorkflow); err != nil {
		log.Fatalln("failed to unmarshal dsl config", err)
	}
	err = validation.ValidateYamlStructure(dslWorkflow)
	if err != nil {
		log.Fatalln("failed to validate dsl config", err)
		panic(err)
	}

	// The client is a heavyweight object that should be created once per process.
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

	workflowOptions := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New(),
		TaskQueue: "dsl",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

}
