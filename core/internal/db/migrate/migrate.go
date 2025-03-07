package migrate

import (
	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/models"
)

func Migrate(db *driver.Driver) error {
	if err := db.Orm.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	if err := db.Orm.AutoMigrate(
		&models.Tenant{},
		&models.Subscription{},
		&models.Workspace{},
		&models.WorkspaceConfig{},
		&models.Upload{},
		&models.Processor{},
		&models.APIKey{},
		&models.Secret{},
	); err != nil {
		return err
	}

	return nil
}
