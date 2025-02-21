package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/uploadpilot/go-core/dsl"
	"github.com/uploadpilot/momentum/internal/activities"
	"github.com/uploadpilot/momentum/internal/config"
	"github.com/uploadpilot/momentum/internal/infra"
	"go.temporal.io/sdk/worker"
)

func startWorker() {
	// Initialize config
	if err := config.Init(); err != nil {
		log.Fatalln("Unable to initialize config.", err)
	}
	config.InitLogger(config.Environment)

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

	// Register all activities
	activities.RegisterActivities(w)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
		os.Exit(1)
	}
}

func main() {
	// Channel to handle OS interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go startWorker()

	// Start a dummy HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Momentum is running")
	})

	server := &http.Server{Addr: ":8085"}
	go func() {
		log.Println("HTTP server started on :8085")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Gracefully shutdown the server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
