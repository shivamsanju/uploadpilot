package main

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := infra.Init(&infra.S3Config{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
		Region:    config.S3Region,
	}); err != nil {
		panic(err)
	}
	sqlDB, err := gorm.Open(postgres.Open(config.PostgresURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}

	err = db.Migrate(sqlDB)
	if err != nil {
		panic(err)
	}
	infra.Log.Info("Migrated database successfully!")

}
