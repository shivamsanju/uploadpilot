package worker

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/uploadpilot/core/internal/workflow/dsl"
	"github.com/uploadpilot/core/internal/workflow/executor"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	lambdaClient   *lambda.Client
	temporalClient client.Client
	taskQueue      string
	wrk            worker.Worker
}

func NewWorker(lambdaClient *lambda.Client, temporalClient client.Client, taskQueue string) *Worker {
	return &Worker{
		lambdaClient:   lambdaClient,
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

func (w *Worker) Start() {
	wrk := worker.New(w.temporalClient, w.taskQueue, worker.Options{
		Identity: w.taskQueue + "-worker",
	})

	wrk.RegisterWorkflow(dsl.SimpleDSLWorkflow)

	exc := executor.NewExecutor(w.lambdaClient)
	wrk.RegisterActivityWithOptions(exc.ExecuteLambdaContainerActivity, activity.RegisterOptions{
		Name: "Executor",
	})

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
