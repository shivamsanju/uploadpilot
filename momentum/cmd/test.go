package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uploadpilot/uploadpilot/momentum/internal/config"
	"github.com/uploadpilot/uploadpilot/momentum/internal/infra"
	"go.temporal.io/api/enums/v1"
)

func mains() {
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

	// w := worker.New(infra.TemporalClient, "dsl", worker.Options{})
	// w.RegisterWorkflow(dsl.SimpleDSLWorkflow)

	// // Register all activities
	// activities.RegisterActivities(w)

	// err = w.Run(worker.InterruptCh())
	// if err != nil {
	// 	log.Fatalln("Unable to start worker", err)
	// 	os.Exit(1)
	// }

	WID := "b381a5b3-292d-472c-b067-0c7d8e2e3599"
	RID := "01950984-d474-7300-8e0a-2ba9f6b0e3a1"
	iter := infra.TemporalClient.GetWorkflowHistory(context.Background(), WID, RID, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	for iter.HasNext() {
		startTime := time.Now() // Capture start time
		event, err := iter.Next()
		if err != nil {
			panic(err)
		}
		fmt.Println(event.GetAttributes())

		eventType := event.GetEventType()
		executionTime := time.Since(startTime)

		switch eventType {
		case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			attributes := event.GetActivityTaskScheduledEventAttributes()
			fmt.Printf("%s | %v | %v - Scheduled | %v\n\n\n",
				ConvertSecondsToTime(event.EventTime.Seconds, event.EventTime.Nanos), eventType, attributes.ActivityType, executionTime)

		case enums.EVENT_TYPE_ACTIVITY_TASK_STARTED:
			attributes := event.GetActivityTaskStartedEventAttributes()
			fmt.Printf("%s | %v | %v - Started | %v\n\n\n",
				ConvertSecondsToTime(event.EventTime.Seconds, event.EventTime.Nanos), eventType, attributes.Attempt, executionTime)

		case enums.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
			attributes := event.GetActivityTaskCompletedEventAttributes()

			fmt.Printf("%s | %v | %v - Success | %v\n\n\n",
				ConvertSecondsToTime(event.EventTime.Seconds, event.EventTime.Nanos), eventType, attributes.Result, executionTime)

		case enums.EVENT_TYPE_ACTIVITY_TASK_FAILED:
			attributes := event.GetActivityTaskFailedEventAttributes()
			fmt.Printf("%s | %v | %v - Failure | %v\n\n\n",
				ConvertSecondsToTime(event.EventTime.Seconds, event.EventTime.Nanos), eventType, attributes.Failure, executionTime)

		default:
			fmt.Printf("%s | %v | %v | %v\n\n\n",
				event.EventTime.AsTime(), eventType, event.String(), event.EventTime.Seconds)
		}

	}
}

func ConvertSecondsToTime(seconds int64, nanos int32) string {
	// Create a time object from Unix seconds and nanoseconds
	t := time.Unix(seconds, int64(nanos)).UTC()
	// Format the time in a readable format
	return t.Format("2006-01-02 15:04:05.000 UTC")
}
