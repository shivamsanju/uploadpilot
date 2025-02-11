package db

import (
	"context"
	"time"

	"github.com/uploadpilot/uploadpilot/common/pkg/cache"
	"github.com/uploadpilot/uploadpilot/common/pkg/db/dbutils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"gorm.io/gorm"
)

type WorkspaceRepo struct {
	wsUserRepo *WorkspaceUserRepo
}

func NewWorkspaceRepo() *WorkspaceRepo {
	return &WorkspaceRepo{
		wsUserRepo: NewWorkspaceUserRepo(),
	}
}

func (wr *WorkspaceRepo) GetAll(ctx context.Context, userID string) ([]models.WorkspaceNameID, error) {
	dbFetchFn := func(workspaces *[]models.WorkspaceNameID) error {
		return sqlDB.WithContext(ctx).Model(&models.Workspace{}).
			Select("id, name").
			Joins("JOIN user_workspaces uw ON workspaces.id = uw.workspace_id").
			Where("uw.user_id = ?", userID).
			Order("workspaces.updated_at desc").
			Find(workspaces).Error
	}

	var workspaces []models.WorkspaceNameID
	cl := cache.NewClient[*[]models.WorkspaceNameID](0)
	key := UserWorkspacesKey(userID)
	if err := cl.Query(ctx, key, &workspaces, dbFetchFn); err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (wr *WorkspaceRepo) Get(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	dbFetchFn := func(workspace *models.Workspace) error {
		return sqlDB.WithContext(ctx).Where("id = ?", workspaceID).
			Select("id, name", "description", "tags").
			First(&workspace).Error
	}

	var workspace models.Workspace
	key := WorkspaceKey(workspaceID)
	cl := cache.NewClient[*models.Workspace](0)
	if err := cl.Query(ctx, key, &workspace, dbFetchFn); err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (wr *WorkspaceRepo) Create(ctx context.Context, workspace *models.Workspace, userWorkspace *models.UserWorkspace, config *models.UploaderConfig) error {
	dbMutateFn := func(workspace *models.Workspace) error {
		return sqlDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

	key := WorkspaceKey(workspace.ID)
	invKeys := []string{UserWorkspacesKey(userWorkspace.UserID)}
	cl := cache.NewClient[*models.Workspace](0)
	if err := cl.Mutate(ctx, key, invKeys, workspace, dbMutateFn, 0); err != nil {
		return err
	}

	return nil
}

func (wr *WorkspaceRepo) Delete(ctx context.Context, workspaceID string) error {
	dbMutateFn := func(workspace *models.Workspace) error {
		return sqlDB.WithContext(ctx).Delete(workspace, "id = ?", workspaceID).Error
	}

	users, err := wr.wsUserRepo.GetUsersInWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}

	key := WorkspaceKey(workspaceID)
	invKeys := []string{
		UserWorkspacesKey(workspaceID),
		WorkspaceConfigKey(workspaceID),
		WorkspaceUsersKey(workspaceID),
	}
	for _, user := range users {
		invKeys = append(invKeys, UserWorkspacesKey(user.ID))
	}

	var workspace models.Workspace
	cl := cache.NewClient[*models.Workspace](0)

	if err := cl.Mutate(ctx, key, []string{}, &workspace, dbMutateFn, 0); err != nil {
		return err
	}

	return nil
}

func (u *WorkspaceRepo) IsSubscriptionActive(ctx context.Context, workspaceID string) (bool, error) {
	var trialEndsAt time.Time

	if err := sqlDB.WithContext(ctx).
		Table("workspaces").
		Joins("JOIN users u ON workspaces.created_by::uuid = u.id").
		Where("workspaces.id = ?", workspaceID).
		Select("u.trial_ends_at").
		Scan(&trialEndsAt).Error; err != nil {
		return false, dbutils.DBError(err)
	}

	return trialEndsAt.After(time.Now()), nil
}

func UserWorkspacesKey(userID string) string {
	return "workspaces:" + userID
}

func WorkspaceKey(workspaceID string) string {
	return "workspace:" + workspaceID
}
