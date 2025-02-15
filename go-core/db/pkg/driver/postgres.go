package driver

import (
	"log"
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
}

func NewPostgresDriver(postgresURI string, config *DBConfig) (*Driver, error) {
	orm, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifeTime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	log.Println("successfully connected to postgres!")

	return &Driver{
		Orm: orm,
	}, nil
}
