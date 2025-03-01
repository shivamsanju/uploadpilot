package worker

import (
	"log"
	"os"

	"github.com/uploadpilot/agent/internal/activities"
	"github.com/uploadpilot/agent/internal/config"
	"github.com/uploadpilot/agent/internal/infra"
	"github.com/uploadpilot/go-common/workflow/dsl"
	"go.temporal.io/sdk/worker"
)

func StartWorker(taskQueue string) {
	// Initialize config
	if err := config.BuildConfig(); err != nil {
		log.Fatalln("Unable to initialize config.", err)
	}
	appConfig := config.GetAppConfig()
	config.InitLogger(appConfig.Environment)

	w := worker.New(infra.TemporalClient, taskQueue, worker.Options{
		Identity: taskQueue + "-worker",
	})
	w.RegisterWorkflow(dsl.SimpleDSLWorkflow)

	// Register all activities
	ar := activities.NewActivityRegistry(w)
	ar.RegisterActivities(w)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
		os.Exit(1)
	}
}
