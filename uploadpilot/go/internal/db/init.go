package db

import (
	"time"

	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sqlDB *gorm.DB
)

func Init() error {
	db, err := gorm.Open(postgres.Open(config.PostgresURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return err
	}

	sqlDB = db

	_db, err := db.DB()
	if err != nil {
		return err
	}

	_db.SetMaxOpenConns(25)                  // Maximum open connections
	_db.SetMaxIdleConns(10)                  // Maximum idle connections
	_db.SetConnMaxLifetime(30 * time.Minute) // Reuse connections for 30 min
	_db.SetConnMaxIdleTime(5 * time.Minute)  // Idle connections max time

	infra.Log.Info("successfully connected to postgres!")

	return nil
}
