package initializer

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"github.com/uploadpilot/core/config"
	"github.com/uploadpilot/core/internal/auth"
	"github.com/uploadpilot/core/internal/clients"
	"github.com/uploadpilot/core/internal/db/cache"
	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/plugins"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/web"
)

func Initialize() (*http.Server, func(), error) {
	if err := config.LoadConfig("./config", "dev", "env"); err != nil {
		return nil, nil, fmt.Errorf("config initialization failed: %w", err)
	}

	setupLogger(config.AppConfig.Environment)

	// Initialize auth
	if err := initAuth(config.AppConfig); err != nil {
		return nil, nil, fmt.Errorf("auth initialization failed: %w", err)
	}

	// Initialize clients
	clients, err := initClients(config.AppConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("clients initialization failed: %w", err)
	}

	// Initialize database
	pgDriver, dbCloseFunc, err := initDatabase(config.AppConfig, clients.RedisClient)
	if err != nil {
		cleanupFunc := getCleanupFunc(clients, nil)
		return nil, cleanupFunc, fmt.Errorf("database initialization failed: %w", err)
	}

	cleanupFunc := getCleanupFunc(clients, dbCloseFunc)

	// Initialize rbac
	accessManager, err := initRBAC(pgDriver, "access_policy")
	if err != nil {
		return nil, nil, fmt.Errorf("rbac initialization failed: %w", err)
	}

	// Initialize repositories
	repos := repo.NewRepositories(pgDriver)

	// Initialize services
	services := services.NewServices(repos, clients, accessManager)

	// Initialize the web server.
	srv, err := web.NewWebserver(config.AppConfig, services)
	if err != nil {
		return nil, nil, fmt.Errorf("web server initialization failed: %w", err)
	}

	return srv, cleanupFunc, nil
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

	s3Opts := &clients.S3Opts{
		AccessKey: appConfig.S3AccessKey,
		SecretKey: appConfig.S3SecretKey,
		Region:    appConfig.S3Region,
	}

	kmsOpts := &clients.KMSOpts{EncryptionKey: appConfig.ApiKeyEncryptionKey}

	return clients.NewAppClients(&clients.ClientOpts{
		RedisOpts:    redisOpts,
		TemporalOpts: temporalOpts,
		S3Opts:       s3Opts,
		KMSOpts:      kmsOpts,
	})
}

func setupLogger(env string) {
	if env == "" {
		env = "development"
	}
	if env == "development" && log.IsTerminal(os.Stderr.Fd()) {
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

func getCleanupFunc(clients *clients.Clients, dbCloseFunc func()) func() {
	return func() {
		_ = clients.RedisClient.Close()
		clients.TemporalClient.Close()
		if dbCloseFunc != nil {
			dbCloseFunc()
		}
	}
}

func initRBAC(dbDriver *driver.Driver, tableName string) (*rbac.AccessManager, error) {
	gormAdapter, err := rbac.NewGormAdapter(dbDriver.Orm, tableName)
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
