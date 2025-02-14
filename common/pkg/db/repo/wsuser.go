package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type WorkspaceUserRepo struct {
	db *db.DB
}

func NewWorkspaceUserRepo(db *db.DB) *WorkspaceUserRepo {
	return &WorkspaceUserRepo{
		db: db,
	}
}

// User methods
func (r *WorkspaceUserRepo) GetUsersInWorkspace(ctx context.Context, workspaceID string) ([]models.WorkspaceUserDetails, error) {
	var users []models.WorkspaceUserDetails
	if err := r.db.Orm.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ?", workspaceID).
		Joins("JOIN users u ON user_workspaces.user_id = u.id").
		Select("user_workspaces.role", "u.name", "u.email", "u.id").
		Find(&users).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return users, nil
}

func (r *WorkspaceUserRepo) AddUserToWorkspace(ctx context.Context, user *models.UserWorkspace) error {
	if err := r.db.Orm.WithContext(ctx).Create(&user).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *WorkspaceUserRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.UserWorkspace{}, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *WorkspaceUserRepo) UpdateUserInWorkspace(ctx context.Context, user *models.UserWorkspace) error {
	if err := r.db.Orm.WithContext(ctx).Save(&user).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *WorkspaceUserRepo) CheckIfUserExistsInWorkspace(ctx context.Context, workspaceID string, userID string) (bool, error) {
	var count int64
	err := r.db.Orm.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, dbutils.DBError(err)
}

func (r *WorkspaceUserRepo) IsOwner(ctx context.Context, workspaceID string, userID string) (bool, error) {
	var count int64
	err := r.db.Orm.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ? AND role = ?", workspaceID, userID, models.UserRoleOwner).
		Count(&count).Error
	return count > 0, dbutils.DBError(err)
}
