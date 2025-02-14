package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type WorkspaceConfigRepo struct {
	db *db.DB
}

func NewWorkspaceConfigRepo(db *db.DB) *WorkspaceConfigRepo {
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
