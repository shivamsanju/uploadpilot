package repo

import (
	"context"

	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/models"
	dbutils "github.com/uploadpilot/core/internal/db/utils"
	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	db *driver.Driver
}

func NewWorkspaceRepo(db *driver.Driver) *WorkspaceRepo {
	return &WorkspaceRepo{
		db: db,
	}
}

func (r *WorkspaceRepo) Get(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.Orm.WithContext(ctx).Where("id = ?", workspaceID).
		Select("id, name", "description", "tags", "created_at").
		Order("created_at DESC").
		First(&workspace).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &workspace, nil
}

func (r *WorkspaceRepo) GetAll(ctx context.Context, tenantID string) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.Orm.WithContext(ctx).Where("tenant_id = ?", tenantID).
		Select("id, name", "description", "tags", "created_at").
		Order("created_at DESC").
		Find(&workspaces).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepo) Create(ctx context.Context, workspace *models.Workspace) error {
	return r.db.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(workspace).Error; err != nil {
			return dbutils.DBError(ctx, r.db.Orm.Logger, err)
		}

		config := models.DefaultWorkspaceConfig
		config.WorkspaceID = workspace.ID
		if err := tx.WithContext(ctx).Create(&config).Error; err != nil {
			return dbutils.DBError(ctx, r.db.Orm.Logger, err)
		}

		return nil
	})
}

func (r *WorkspaceRepo) Delete(ctx context.Context, workspaceID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", workspaceID).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}
