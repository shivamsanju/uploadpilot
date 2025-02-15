package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/driver"
	cacheplugins "github.com/uploadpilot/uploadpilot/go-core/db/pkg/plugins/cache"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/uploadpilot/manager/internal/auth"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/uploadpilot/manager/web"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	srv, err := initialize()
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
			log.Fatal(fmt.Errorf("graceful server shutdown failed: %w", err))
			return
		}

		log.Println("server gracefully stopped")
	}(wg)

	// Start the web server.
	log.Printf("starting web server on port %d\n", config.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Errorf("server initialization failed: %w", err))
	}
	wg.Wait()
	log.Println("server exited")

}

func initialize() (*http.Server, error) {
	if err := config.Init(); err != nil {
		return nil, fmt.Errorf("config initialization failed: %w", err)
	}

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
	if config.RedisTLS {
		redisOpts.TLSConfig = &tls.Config{}
	}

	err := infra.Init(&infra.InfraOpts{
		S3Opts:    s3Opts,
		RedisOpts: redisOpts,
	})

	if err != nil {
		return nil, fmt.Errorf("infra initialization failed: %w", err)
	}

	// Initialize authentication
	if err := auth.InitSessionStore(); err != nil {
		return nil, fmt.Errorf("auth initialization failed: %w", err)
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

	// Initialize the web server.
	repos := repo.NewRepositories(pgDriver)
	svcs := svc.NewServices(repos)
	srv, err := web.InitWebServer(svcs)
	if err != nil {
		return nil, fmt.Errorf("web server initialization failed: %w", err)
	}

	return srv, nil
}

func cleanup() {
	// Perform any cleanup here, e.g., closing database connections, stopping services.
	// Example: if err := db.Close(); err != nil { return err }
	log.Println("performing cleanup...")
}
