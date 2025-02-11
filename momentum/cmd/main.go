package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/kms"
	"github.com/uploadpilot/uploadpilot/momentum/internal/config"
	"github.com/uploadpilot/uploadpilot/momentum/web"
)

func mains() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	srv, err := initServices()
	if err != nil {
		cleanup()
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sig := <-sigChan
		log.Printf("received shutdown signal: %s\n", sig)

		cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(wrapError("graceful server shutdown failed", err))
			return
		}

		log.Println("server gracefully stopped")
	}(wg)

	// Start the web server.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(wrapError("server initialization failed", err))
	}

	wg.Wait()
	log.Println("server exited")

}

func initServices() (*http.Server, error) {
	// Initialize configuration.
	if err := config.Init(); err != nil {
		return nil, wrapError("config initialization failed", err)
	}

	// Initialize KMS
	if err := kms.Init(string(config.EncryptionKey)); err != nil {
		return nil, wrapError("kvm initialization failed", err)
	}

	// Initialize infra.
	if err := infra.Init(&infra.S3Config{
		AccessKey: "",
		SecretKey: "",
		Region:    "",
	}); err != nil {
		return nil, wrapError("infra initialization failed", err)
	}

	// Initialize database.
	if err := db.Init(config.PostgresURI); err != nil {
		return nil, wrapError("database initialization failed", err)
	}

	// Initialize the web server.
	srv, err := web.Init()
	if err != nil {
		return nil, wrapError("web server initialization failed", err)
	}

	return srv, nil
}

func cleanup() {
	// Perform any cleanup here, e.g., closing database connections, stopping services.
	// Example: if err := db.Close(); err != nil { return err }
	log.Println("performing cleanup...")
}

// wrapError provides better error context.
func wrapError(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}
