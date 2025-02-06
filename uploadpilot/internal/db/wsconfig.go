package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/cache"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
)

type WorkspaceConfigRepo struct {
}

func NewWorkspaceConfigRepo() *WorkspaceConfigRepo {
	return &WorkspaceConfigRepo{}
}

// Config methods
func (wr *WorkspaceConfigRepo) GetConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	key := WorkspaceConfigKey(workspaceID)
	dbFetchFn := func(config *models.UploaderConfig) error {
		return sqlDB.WithContext(ctx).Omit("Workspace").Where("workspace_id = ?", workspaceID).
			First(config).Error
	}
	var config models.UploaderConfig
	cl := cache.NewClient[*models.UploaderConfig](0)
	if err := cl.Query(ctx, key, &config, dbFetchFn); err != nil {
		return nil, err
	}
	return &config, nil
}

func (wr *WorkspaceConfigRepo) SetConfig(ctx context.Context, config *models.UploaderConfig) error {
	key := WorkspaceConfigKey(config.WorkspaceID)

	dbMutateFn := func(config *models.UploaderConfig) error {
		return sqlDB.WithContext(ctx).Save(config).Error
	}
	cl := cache.NewClient[*models.UploaderConfig](0)
	return cl.Mutate(ctx, key, []string{}, config, dbMutateFn, 0)
}

func WorkspaceConfigKey(workspaceID string) string {
	return "workspace:" + workspaceID + ":config"
}
