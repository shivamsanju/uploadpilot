package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/phuslu/log"

	"github.com/uploadpilot/core/config"
	initializer "github.com/uploadpilot/core/init"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	srv, wrk, timeoutMarkerCron, cleanupFunc, err := initializer.Initialize()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize")
		if cleanupFunc != nil {
			cleanupFunc()
		}
		os.Exit(1)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("received signal")

		if cleanupFunc == nil {
			cleanupFunc()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("graceful server shutdown failed")
			return
		}

		log.Info().Msg("server gracefully stopped")
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		wrk.Start()
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		timeoutMarkerCron.Start()
	}(wg)

	// Start the web server.
	log.Info().Int("port", config.AppConfig.Port).Msg("starting web server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("failed to start web server")
	}

	wg.Wait()
	log.Info().Msg("exited")
}
