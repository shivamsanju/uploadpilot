package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/cache"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
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

func (wr *WorkspaceRepo) GetAll(ctx context.Context, userID string) ([]dto.WorkspaceNameID, error) {
	dbFetchFn := func(workspaces *[]dto.WorkspaceNameID) error {
		return sqlDB.WithContext(ctx).Model(&models.Workspace{}).
			Select("id, name").
			Joins("JOIN user_workspaces uw ON workspaces.id = uw.workspace_id").
			Where("uw.user_id = ?", userID).
			Order("workspaces.updated_at ASC").
			Find(workspaces).Error
	}

	var workspaces []dto.WorkspaceNameID
	cl := cache.NewClient[*[]dto.WorkspaceNameID](0)
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
				return utils.DBError(err)
			}

			userWorkspace.WorkspaceID = workspace.ID
			if err := tx.WithContext(ctx).Create(&userWorkspace).Error; err != nil {
				return utils.DBError(err)
			}

			config.WorkspaceID = workspace.ID
			if err := tx.WithContext(ctx).Create(&config).Error; err != nil {
				return utils.DBError(err)
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
		WorkSpaceConfigKey(workspaceID),
		WorkSpaceUsersKey(workspaceID),
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

func UserWorkspacesKey(userID string) string {
	return "workspaces:" + userID
}

func WorkspaceKey(workspaceID string) string {
	return "workspace:" + workspaceID
}
