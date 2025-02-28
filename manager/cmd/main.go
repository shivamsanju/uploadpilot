package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/phuslu/log"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/go-core/common/vault"
	"github.com/uploadpilot/go-core/db/pkg/driver"
	cacheplugins "github.com/uploadpilot/go-core/db/pkg/plugins/cache"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/manager/web"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	srv, err := initialize()
	if err != nil {
		cleanup()
		panic(err)
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
			log.Error().Err(err).Msg("graceful server shutdown failed")
			return
		}

		log.Info().Msg("server gracefully stopped")
	}(wg)

	// Start the web server.
	log.Info().Int("port", config.Port).Msg("starting web server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("failed to start web server")
	}
	wg.Wait()
	log.Info().Msg("exited")

}

func initialize() (*http.Server, error) {
	if err := config.Init(); err != nil {
		return nil, fmt.Errorf("config initialization failed: %w", err)
	}
	config.InitLogger(config.Environment)

	// Initialize infra
	s3Opts := &infra.S3Options{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
		Region:    config.S3Region,
	}

	redisOpts := &redis.Options{
		Addr:     config.RedisAddr,
		Username: config.RedisUsername,
		Password: config.RedisPassword,
	}

	temporalOpts := &infra.TemporalOptions{
		Namespace: config.TemporalNamespace,
		HostPort:  config.TemporalHostPort,
		APIKey:    config.TemporalAPIKey,
	}

	superTokensOpts := &infra.SuperTokensOptions{
		ConnectionURI:   config.SuperTokensEndpoint,
		APIKey:          config.SupertokensAPIKey,
		AppName:         config.AppName,
		APIBasePath:     "/auth",
		WebsiteBasePath: "/auth",
		APIDomain:       config.SelfEndpoint,
		WebsiteDomain:   config.FrontendURI,
		Providers: []infra.Provider{
			{
				Key:          "google",
				ClientID:     config.GoogleClientID,
				ClientSecret: config.GoogleClientSecret,
			},
			{
				Key:          "github",
				ClientID:     config.GithubClientID,
				ClientSecret: config.GithubClientSecret,
			},
		},
	}

	if config.RedisTLS {
		redisOpts.TLSConfig = &tls.Config{}
	}

	err := infra.Init(&infra.InfraOpts{
		S3Opts:          s3Opts,
		RedisOpts:       redisOpts,
		TemporalOpts:    temporalOpts,
		SuperTokensOpts: superTokensOpts,
	})

	if err != nil {
		return nil, fmt.Errorf("infra initialization failed: %w", err)
	}

	// Initialize database
	pgDriver, err := driver.NewPostgresDriver(config.PostgresURI, &driver.DBConfig{
		MaxOpenConn:     10,
		MaxIdleConn:     5,
		ConnMaxLifeTime: time.Minute * 30,
		ConnMaxIdleTime: time.Minute * 5,
	})
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	// Add caching layer
	rcp := cacheplugins.NewRedisCachesPlugin(infra.RedisClient)
	pgDriver.Orm.Use(rcp)

	// Initialize KMS
	kms, err := vault.NewKMS(config.ApiKeyEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("kms initialization failed: %w", err)
	}
	// Initialize the web server.
	repos := repo.NewRepositories(pgDriver)
	svcs := svc.NewServices(repos, kms)
	srv, err := web.InitWebServer(svcs)
	if err != nil {
		return nil, fmt.Errorf("web server initialization failed: %w", err)
	}

	return srv, nil
}

func cleanup() {
	// Perform any cleanup here, e.g., closing database connections, stopping services.
	// Example: if err := db.Close(); err != nil { return err }
	log.Info().Msg("performing cleanup...")
}
