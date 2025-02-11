package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/cache"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type WorkspaceUserRepo struct {
}

func NewWorkspaceUserRepo() *WorkspaceUserRepo {
	return &WorkspaceUserRepo{}
}

// User methods
func (wr *WorkspaceUserRepo) GetUsersInWorkspace(ctx context.Context, workspaceID string) ([]models.WorkspaceUserDetails, error) {
	dbFetchFn := func(users *[]models.WorkspaceUserDetails) error {
		return sqlDB.WithContext(ctx).Model(&models.UserWorkspace{}).
			Where("workspace_id = ?", workspaceID).
			Joins("JOIN users u ON user_workspaces.user_id = u.id").
			Select("user_workspaces.role", "u.name", "u.email", "u.id").
			Find(users).Error
	}
	var users []models.WorkspaceUserDetails

	cl := cache.NewClient[*[]models.WorkspaceUserDetails](0)
	key := WorkspaceUsersKey(workspaceID)
	if err := cl.Query(ctx, key, &users, dbFetchFn); err != nil {
		return nil, err
	}
	return users, nil
}

func (wr *WorkspaceUserRepo) AddUserToWorkspace(ctx context.Context, user *models.UserWorkspace) error {
	mutateDbFn := func(user *models.UserWorkspace) error {
		return sqlDB.WithContext(ctx).Create(&user).Error
	}
	cl := cache.NewClient[*models.UserWorkspace](0)
	key := WorkspaceUserKey(user.WorkspaceID, user.UserID)
	invkeys := []string{UserWorkspacesKey(user.UserID), WorkspaceUsersKey(user.WorkspaceID)}
	if err := cl.Mutate(ctx, key, invkeys, user, mutateDbFn, 0); err != nil {
		return err
	}

	return nil
}

func (wr *WorkspaceUserRepo) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	dbMutateFn := func(user *models.UserWorkspace) error {
		return sqlDB.WithContext(ctx).Delete(&user, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error
	}

	cl := cache.NewClient[*models.UserWorkspace](0)
	key := WorkspaceUserKey(workspaceID, userID)
	invkeys := []string{UserWorkspacesKey(userID), WorkspaceUsersKey(workspaceID)}
	if err := cl.Mutate(ctx, key, invkeys, nil, dbMutateFn, 0); err != nil {
		return err
	}

	return nil
}

func (wr *WorkspaceUserRepo) UpdateUserInWorkspace(ctx context.Context, user *models.UserWorkspace) error {
	dbMutateFn := func(user *models.UserWorkspace) error {
		return sqlDB.WithContext(ctx).Save(&user).Error
	}

	cl := cache.NewClient[*models.UserWorkspace](0)
	key := WorkspaceUserKey(user.WorkspaceID, user.UserID)
	invkeys := []string{UserWorkspacesKey(user.UserID), WorkspaceUsersKey(user.WorkspaceID)}
	if err := cl.Mutate(ctx, key, invkeys, user, dbMutateFn, 0); err != nil {
		return err
	}

	return nil
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

func WorkspaceUsersKey(workspaceID string) string {
	return "workspace:" + workspaceID + ":users"
}

func WorkspaceUserKey(workspaceID string, userID string) string {
	return "workspace:" + workspaceID + ":user:" + userID
}
