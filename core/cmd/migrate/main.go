package main

import (
	"log"
	"os"

	"github.com/uploadpilot/core/config"
	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/migrate"
	"gorm.io/gorm/logger"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Environment not specified. Usage: go run cmd/migrate/main.go [dev|prod|test]")
	}

	env := os.Args[1]
	config.LoadConfig("./config", env, "env")

	pgDriver, closeFunc, err := driver.NewPostgresDriver(config.AppConfig.PostgresURI, &driver.DBConfig{
		LogMode:     logger.Info,
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	})
	if err != nil {
		panic(err)
	}
	defer closeFunc()

	err = migrate.Migrate(pgDriver)
	if err != nil {
		panic(err)
	}

	log.Println("migrated database successfully!")
}
