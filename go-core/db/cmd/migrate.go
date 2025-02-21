package main

import (
	"log"

	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/migrate"
)

func main() {
	dbUri := "postgresql://postgres.wjxdjummbehatmlfrqoa:sanjushivam@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres"
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
