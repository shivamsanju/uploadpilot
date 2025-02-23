package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/phuslu/log"
	"github.com/uploadpilot/uploader/internal/clients"
	"github.com/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploader/internal/service"
	"github.com/uploadpilot/uploader/web"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	if err := config.BuildConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	appConfig := config.GetAppConfig()
	config.InitLogger(appConfig.Environment)
	coreClient := clients.NewCoreServiceClient(appConfig.CoreServiceEndpoint, appConfig.CoreServiceAPIKey)
	svc := service.NewUploadService(coreClient)

	srv, err := web.InitWebserver(svc)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize web server routes")
		os.Exit(1)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("received shutdown signal")

		cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("graceful server shutdown failed")
			os.Exit(1)
		}

		log.Info().Msg("server gracefully stopped")
	}(wg)

	// Start the web server.
	log.Info().Int("port", appConfig.Port).Msg("starting web server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("failed to start web server")
		os.Exit(1)
	}

	wg.Wait()
}

func cleanup() {
	// Perform any cleanup here, e.g., closing database connections, stopping services.
	// Example: if err := db.Close(); err != nil { return err }
	log.Info().Msg("performing cleanup...")
}
