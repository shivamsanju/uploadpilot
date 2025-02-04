package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type WorkspaceUserRepo struct {
}

func NewWorkspaceUserRepo() *WorkspaceUserRepo {
	return &WorkspaceUserRepo{}
}

// User methods
func (wr *WorkspaceUserRepo) GetUsersInWorkspace(ctx context.Context, workspaceID string) ([]dto.WorkspaceUser, error) {

	var users []dto.WorkspaceUser
	err := sqlDB.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ?", workspaceID).
		Joins("JOIN users u ON user_workspaces.user_id = u.id").
		Select("user_workspaces.role", "u.name", "u.email", "u.id").
		Find(&users).Error

	if err != nil {
		return nil, utils.DBError(err)
	}
	return users, nil
}

func (wr *WorkspaceUserRepo) AddUserToWorkspace(ctx context.Context, user *models.UserWorkspace) ([]dto.WorkspaceUser, error) {
	if err := sqlDB.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, utils.DBError(err)
	}

	return wr.GetUsersInWorkspace(ctx, user.WorkspaceID)
}

func (wr *WorkspaceUserRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) ([]dto.WorkspaceUser, error) {
	if err := sqlDB.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Delete(&models.UserWorkspace{}).Error; err != nil {
		return nil, utils.DBError(err)
	}

	return wr.GetUsersInWorkspace(ctx, workspaceID)
}

func (wr *WorkspaceUserRepo) UpdateUserInWorkspace(ctx context.Context, user *models.UserWorkspace) ([]dto.WorkspaceUser, error) {
	if err := sqlDB.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, utils.DBError(err)
	}

	return wr.GetUsersInWorkspace(ctx, user.WorkspaceID)
}

func (wr *WorkspaceUserRepo) CheckIfUserExistsInWorkspace(ctx context.Context, workspaceID string, userID string) (bool, error) {
	var count int64
	err := sqlDB.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&count).Error
	return count > 0, err
}

func (wr *WorkspaceUserRepo) IsOwner(ctx context.Context, workspaceID string, userID string) (bool, error) {
	var count int64
	err := sqlDB.WithContext(ctx).Model(&models.UserWorkspace{}).
		Where("workspace_id = ? AND user_id = ? AND role = ?", workspaceID, userID, models.UserRoleOwner).
		Count(&count).Error
	return count > 0, err
}

func WorkSpaceUsersKey(workspaceID string) string {
	return "workspace:" + workspaceID + ":users"
}

func WorkSpaceUserKey(workspaceID string, userID string) string {
	return "workspace:" + workspaceID + ":user:" + userID
}
