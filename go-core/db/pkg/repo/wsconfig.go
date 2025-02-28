package repo

import (
	"context"

	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/models"
	dbutils "github.com/uploadpilot/go-core/db/pkg/utils"
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
func (r *WorkspaceConfigRepo) GetConfig(ctx context.Context, workspaceID string) (*models.WorkspaceConfig, error) {
	var config models.WorkspaceConfig
	if err := r.db.Orm.WithContext(ctx).Omit("Workspace").Where("workspace_id = ?", workspaceID).
		First(&config).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &config, nil
}

func (r *WorkspaceConfigRepo) SetConfig(ctx context.Context, config *models.WorkspaceConfig) error {
	if err := r.db.Orm.WithContext(ctx).Save(config).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
