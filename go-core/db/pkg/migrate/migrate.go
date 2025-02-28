package migrate

import (
	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/models"
)

func Migrate(db *driver.Driver) error {
	if err := db.Orm.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	if err := db.Orm.Exec(`
	DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'allowed_sources') THEN
		CREATE TYPE allowed_sources AS ENUM (
			'FileUpload',
			'Audio',
			'Webcamera',
			'ScreenCapture',
			'Box',
			'Dropbox',
			'Facebook',
			'GoogleDrive',
			'GooglePhotos',
			'Instagram',
			'OneDrive',
			'Unsplash',
			'Url',
			'Zoom'
		);
	    END IF;
END $$;
	`).Error; err != nil {
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
		&models.APIKeyPermission{},
		&models.Secret{},
	); err != nil {
		return err
	}

	return nil
}
