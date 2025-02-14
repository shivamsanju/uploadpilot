package main

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/db/migrate"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}
	if err := infra.Init(nil); err != nil {
		panic(err)
	}

	sqlDb, err := db.NewPostgresDB(config.PostgresURI, &db.DBConfig{})
	if err != nil {
		panic(err)
	}

	err = migrate.Migrate(sqlDb)
	if err != nil {
		panic(err)
	}
	infra.Log.Info("Migrated database successfully!")

}
