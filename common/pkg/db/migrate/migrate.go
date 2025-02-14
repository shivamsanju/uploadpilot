package migrate

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

func Migrate(db *db.DB) error {
	if err := db.Orm.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	if err := db.Orm.Exec(`
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'upload_log_level') THEN
        CREATE TYPE upload_log_level AS ENUM ('info', 'warn', 'error');
    END IF;
END $$;
`).Error; err != nil {
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
		&models.User{},
		&models.Workspace{},
		&models.UserWorkspace{},
		&models.UploaderConfig{},
		&models.Upload{},
		&models.Processor{},
		&models.UploadLog{},
	); err != nil {
		return err
	}

	return nil
}
