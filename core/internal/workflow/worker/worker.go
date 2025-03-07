package worker

import (
	"log"
	"os"

	"github.com/uploadpilot/core/pkg/dsl"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	temporalClient client.Client
	taskQueue      string
	wrk            worker.Worker
}

func NewWorker(temporalClient client.Client, taskQueue string) *Worker {
	return &Worker{
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

func (w *Worker) Start() {
	wrk := worker.New(w.temporalClient, w.taskQueue, worker.Options{
		Identity: w.taskQueue + "-worker",
	})

	wrk.RegisterWorkflow(dsl.SimpleDSLWorkflow)

	ar := NewActivityRegistry(wrk)
	ar.RegisterActivities(wrk)

	w.wrk = wrk
	err := wrk.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
		os.Exit(1)
	}

}

func (w *Worker) Stop() {
	if w.wrk != nil {
		w.wrk.Stop()
	}
}
