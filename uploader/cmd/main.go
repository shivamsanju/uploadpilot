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
	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/listeners"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc"
	"github.com/uploadpilot/uploadpilot/uploader/web"
)

func main() {
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

	// Initialize infra.
	s3Opts := &infra.S3Options{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
		Region:    config.S3Region,
	}

	temporalOpts := &infra.TemporalOptions{
		Namespace: config.TemporalNamespace,
		HostPort:  config.TemporalHostPort,
		APIKey:    config.TemporalAPIKey,
	}

	if err := infra.Init(&infra.InfraOpts{
		S3Opts:       s3Opts,
		TemporalOpts: temporalOpts,
	}); err != nil {
		return nil, wrapError("infra initialization failed", err)
	}

	// Initialize database.
	db, err := db.NewPostgresDB(config.PostgresURI, &db.DBConfig{
		MaxOpenConn:     10,
		MaxIdleConn:     5,
		ConnMaxLifeTime: time.Minute * 30,
		ConnMaxIdleTime: time.Minute * 5,
	})
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	// Initialize the services
	repos := repo.NewRepositories(db)
	svcs := svc.NewServices(repos)

	// Initialize listeners.
	listeners.StartListeners(svcs)

	// Initialize the uploader web server.
	srv, err := web.InitWebserver(svcs)
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
