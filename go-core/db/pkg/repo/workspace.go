package repo

import (
	"context"
	"time"

	"github.com/uploadpilot/go-core/db/pkg/driver"
	"github.com/uploadpilot/go-core/db/pkg/models"
	dbutils "github.com/uploadpilot/go-core/db/pkg/utils"
	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	wsUserRepo *WorkspaceUserRepo
	db         *driver.Driver
}

func NewWorkspaceRepo(db *driver.Driver) *WorkspaceRepo {
	return &WorkspaceRepo{
		db:         db,
		wsUserRepo: NewWorkspaceUserRepo(db),
	}
}

func (r *WorkspaceRepo) GetAll(ctx context.Context, userID string) ([]models.WorkspaceNameID, error) {
	var workspaces []models.WorkspaceNameID
	if err := r.db.Orm.WithContext(ctx).Model(&models.Workspace{}).
		Select("id, name").
		Joins("JOIN user_workspaces uw ON workspaces.id = uw.workspace_id").
		Where("uw.user_id = ?", userID).
		Order("workspaces.updated_at desc").
		Find(&workspaces).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepo) Get(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.Orm.WithContext(ctx).Where("id = ?", workspaceID).
		Select("id, name", "description", "tags").
		First(&workspace).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return &workspace, nil
}

func (r *WorkspaceRepo) Create(ctx context.Context, workspace *models.Workspace, userWorkspace *models.UserWorkspace, config *models.UploaderConfig) error {
	return r.db.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(workspace).Error; err != nil {
			return dbutils.DBError(err)
		}

		userWorkspace.WorkspaceID = workspace.ID
		if err := tx.WithContext(ctx).Create(&userWorkspace).Error; err != nil {
			return dbutils.DBError(err)
		}

		config.WorkspaceID = workspace.ID
		if err := tx.WithContext(ctx).Create(&config).Error; err != nil {
			return dbutils.DBError(err)
		}

		return nil
	})
}

func (r *WorkspaceRepo) Delete(ctx context.Context, workspaceID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", workspaceID).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *WorkspaceRepo) GetSubscription(ctx context.Context, workspaceID string) (bool, time.Time, error) {
	var trialEndsAt time.Time

	if err := r.db.Orm.WithContext(ctx).
		Table("workspaces").
		Joins("JOIN users u ON workspaces.created_by::uuid = u.id").
		Where("workspaces.id = ?", workspaceID).
		Select("u.trial_ends_at").
		Scan(&trialEndsAt).Error; err != nil {
		return false, time.Time{}, dbutils.DBError(err)
	}

	return trialEndsAt.After(time.Now()), trialEndsAt, nil
}

func (r *WorkspaceRepo) GetApiKeySalt(ctx context.Context, workspaceID string) (string, error) {
	var apiKeySalt string
	if err := r.db.Orm.WithContext(ctx).
		Table("workspaces").
		Where("id = ?", workspaceID).
		Select("api_key_salt").
		Scan(&apiKeySalt).Error; err != nil {
		return "", dbutils.DBError(err)
	}

	return apiKeySalt, nil
}
