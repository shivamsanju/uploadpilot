package driver

import (
	"time"

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
