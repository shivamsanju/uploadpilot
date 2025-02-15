package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	dbutils "github.com/uploadpilot/uploadpilot/go-core/db/pkg/utils"
)

type WorkspaceConfigRepo struct {
	db *driver.Driver
}

func NewWorkspaceConfigRepo(db *driver.Driver) *WorkspaceConfigRepo {
	return &WorkspaceConfigRepo{
		db: db,
	}
}

// Config methods
func (r *WorkspaceConfigRepo) GetConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	var config models.UploaderConfig
	if err := r.db.Orm.WithContext(ctx).Omit("Workspace").Where("workspace_id = ?", workspaceID).
		First(&config).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return &config, nil
}

func (r *WorkspaceConfigRepo) SetConfig(ctx context.Context, config *models.UploaderConfig) error {
	if err := r.db.Orm.WithContext(ctx).Save(config).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}
