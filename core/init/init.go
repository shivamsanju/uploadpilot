package initializer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/uploadpilot/core/config"
	"github.com/uploadpilot/core/internal/auth"
	"github.com/uploadpilot/core/internal/clients"
	"github.com/uploadpilot/core/internal/db/cache"
	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/plugins"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/internal/workflow"
	"github.com/uploadpilot/core/web"
)

func Initialize() (*http.Server, *workflow.Worker, *cron.Cron, func(), error) {
	environment := GetEnvironment()
	setupLogger(environment)

	if err := config.LoadConfig("./config", environment, "env"); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("config initialization failed: %w", err)
	}

	// Initialize auth
	if err := initAuth(config.AppConfig); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("auth initialization failed: %w", err)
	}

	// Initialize clients
	clients, err := initClients(config.AppConfig)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("clients initialization failed: %w", err)
	}

	// Initialize database
	pgDriver, dbCloseFunc, err := initDatabase(config.AppConfig, clients.RedisClient)
	if err != nil {
		cleanupFunc := getCleanupFunc(clients, nil, nil)
		return nil, nil, nil, cleanupFunc, fmt.Errorf("database initialization failed: %w", err)
	}

	// Initialize rbac
	accessManager, err := initRBAC(config.AppConfig, environment)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("rbac initialization failed: %w", err)
	}

	// Initialize repositories
	repos := repo.NewRepositories(pgDriver)

	// Initialize services
	services := services.NewServices(repos, clients, accessManager)

	// Initialize worker
	wrk := workflow.NewWorker(clients.LambdaClient, clients.TemporalClient, config.AppConfig.WorkerTaskQueue)

	// Initialize cron to mark timed out uploads
	timeoutMarkerCron := NewMarkTimedOutUploadsRoutine(repos.UploadRepo)

	// Initialize the web server.
	srv, err := web.NewWebserver(config.AppConfig, services)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("web server initialization failed: %w", err)
	}

	cleanupFunc := getCleanupFunc(clients, dbCloseFunc, wrk)

	return srv, wrk, timeoutMarkerCron, cleanupFunc, nil
}

func initDatabase(appConfig *config.Config, redisClient *redis.Client) (*driver.Driver, func(), error) {
	if appConfig.UseCache && redisClient == nil {
		return nil, nil, fmt.Errorf("redis client must be provided when using cache")
	}

	pgDriver, closeFunc, err := driver.NewPostgresDriver(appConfig.PostgresURI, &driver.DBConfig{
		MaxOpenConn:     10,
		MaxIdleConn:     5,
		ConnMaxLifeTime: time.Minute * 30,
		ConnMaxIdleTime: time.Minute * 5,
	})
	if err != nil {
		return nil, nil, err
	}

	if appConfig.UseCache {
		cacher := cache.NewRedisCacher(redisClient)
		rcp := plugins.NewCachePlugin(cacher)
		pgDriver.Orm.Use(rcp)
	}

	return pgDriver, closeFunc, nil
}

func initAuth(appConfig *config.Config) error {
	superTokensOpts := &auth.SuperTokensOptions{
		ConnectionURI:   appConfig.SuperTokensEndpoint,
		APIKey:          appConfig.SupertokensAPIKey,
		AppName:         appConfig.AppName,
		APIBasePath:     appConfig.APIBaseAuthPath,
		WebsiteBasePath: appConfig.WebsiteBaseAuthPath,
		APIDomain:       appConfig.SelfEndpoint,
		WebsiteDomain:   appConfig.FrontendURI,
		Providers: []auth.Provider{
			{
				Key:          "google",
				ClientID:     appConfig.GoogleClientID,
				ClientSecret: appConfig.GoogleClientSecret,
			},
			{
				Key:          "github",
				ClientID:     appConfig.GithubClientID,
				ClientSecret: appConfig.GithubClientSecret,
			},
		},
	}

	return auth.InitSuperTokens(superTokensOpts)
}

func initClients(appConfig *config.Config) (*clients.Clients, error) {
	redisOpts := &clients.RedisOpts{
		Addr:     appConfig.RedisAddr,
		Username: appConfig.RedisUsername,
		Password: appConfig.RedisPassword,
		TLS:      appConfig.RedisTLS,
	}

	temporalOpts := &clients.TemporalOpts{
		Namespace: appConfig.TemporalNamespace,
		HostPort:  appConfig.TemporalHostPort,
		APIKey:    appConfig.TemporalAPIKey,
	}

	awsOpts := &clients.AwsOpts{
		AccessKey: appConfig.S3AccessKey,
		SecretKey: appConfig.S3SecretKey,
		Region:    appConfig.S3Region,
	}

	kmsOpts := &clients.KMSOpts{EncryptionKey: appConfig.ApiKeyEncryptionKey}

	return clients.NewAppClients(&clients.ClientOpts{
		RedisOpts:    redisOpts,
		TemporalOpts: temporalOpts,
		S3Opts:       awsOpts,
		LambdaOpts:   awsOpts,
		KMSOpts:      kmsOpts,
	})
}

func setupLogger(env string) {
	if env == "dev" && log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}
}

func getCleanupFunc(clients *clients.Clients, dbCloseFunc func(), wrk *workflow.Worker) func() {
	return func() {
		_ = clients.RedisClient.Close()
		clients.TemporalClient.Close()
		if dbCloseFunc != nil {
			dbCloseFunc()
		}
		if wrk != nil {
			wrk.Stop()
		}
	}
}

func initRBAC(config *config.Config, env string) (*rbac.AccessManager, error) {
	gormAdapter, err := rbac.NewPgAdapter(config.PostgresURI, env)
	if err != nil {
		return nil, err
	}

	manager, err := rbac.NewAccessManager(gormAdapter)
	if err != nil {
		return nil, err
	}

	if err := manager.SetupPolicy(); err != nil {
		return nil, err
	}

	return manager, nil
}

func GetEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		env = "dev"
	} else if env == "production" {
		env = "prod"
	} else if env == "testing" {
		env = "test"
	}

	if env == "" {
		env = "dev"
	}

	os.Setenv("ENVIRONMENT", env)
	return env
}

func NewMarkTimedOutUploadsRoutine(uploadRepo *repo.UploadRepo) *cron.Cron {
	c := cron.New()

	// Every 5 minutes
	c.AddFunc("*/5 * * * *", func() {
		err := uploadRepo.BulkMarkTimedOut(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to mark timed out uploads")
		}
		log.Debug().Msg("marked timed out uploads")
	})

	return c
}
