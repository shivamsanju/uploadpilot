package main

import (
	"log"
	"os"

	"github.com/uploadpilot/uploadpilot/common/pkg/dsl"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/momentum/internal/activities"
	"github.com/uploadpilot/uploadpilot/momentum/internal/config"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	if err := config.Init(); err != nil {
		log.Fatalln("Unable to initialize config.", err)
	}

	err := infra.Init(&infra.InfraOpts{
		TemporalOpts: &infra.TemporalOptions{
			Namespace: config.TemporalNamespace,
			HostPort:  config.TemporalHostPort,
			APIKey:    config.TemporalAPIKey,
		},
	})

	if err != nil {
		log.Fatalln("Unable to initialize infra.", err)
		os.Exit(1)
	}
	defer infra.TemporalClient.Close()

	w := worker.New(infra.TemporalClient, "dsl", worker.Options{})
	w.RegisterWorkflow(dsl.SimpleDSLWorkflow)
	w.RegisterActivity(&activities.SampleActivities{})

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
		os.Exit(1)
	}
}
