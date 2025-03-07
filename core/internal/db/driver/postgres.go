package driver

import (
	"fmt"
	"regexp"
	"time"

	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	MaxOpenConn     int           // Maximum open connections
	MaxIdleConn     int           // Maximum idle connections
	ConnMaxLifeTime time.Duration // Reuse connections for ConnMaxLifeTime
	ConnMaxIdleTime time.Duration // Idle connections max time
	LogMode         logger.LogLevel
}

func NewPostgresDriver(postgresURI string, config *DBConfig) (*Driver, func(), error) {
	if err := ensureDatabase(postgresURI); err != nil {
		log.Error().Err(err).Msg("failed to ensure database")
		return nil, nil, err
	}

	if config.LogMode == 0 {
		config.LogMode = logger.Warn
	}

	// Initialize Gorm with Phulsu Log
	orm, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogMode),
	})

	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifeTime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	closeFunc := func() {
		sqlDB.Close()
	}

	return &Driver{
		Orm: orm,
	}, closeFunc, nil
}

func ensureDatabase(postgresURI string) error {
	re := regexp.MustCompile(`(.*/)([^/]+)$`)
	matches := re.FindStringSubmatch(postgresURI)
	if len(matches) < 3 {
		return fmt.Errorf("invalid PostgreSQL URI format")
	}

	baseURI := matches[1]
	dbName := matches[2]

	adminURI := baseURI + "postgres" // Connect to the default 'postgres' database
	g, err := gorm.Open(postgres.Open(adminURI))
	if err != nil {
		return fmt.Errorf("failed to connect to admin database: %w", err)
	}
	db, err := g.DB()
	if err != nil {
		return fmt.Errorf("failed to get admin database connection: %w", err)
	}
	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_database WHERE datname = '%s')", dbName)
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE \"%s\"", dbName)
		_, err := db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Info().Str("database", dbName).Msg("database created")
	}

	return nil
}
