package main

import (
	"log"
	"os"

	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/migrate"
)

func main() {
	dbUri := os.Getenv("POSTGRES_URI")
	pgDriver, err := driver.NewPostgresDriver(dbUri, &driver.DBConfig{})
	if err != nil {
		panic(err)
	}

	err = migrate.Migrate(pgDriver)
	if err != nil {
		panic(err)
	}
	log.Println("migrated database successfully!")

}
