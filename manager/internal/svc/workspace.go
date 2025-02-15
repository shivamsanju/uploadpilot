package svc

import (
	"context"
	"errors"
	"fmt"

	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

var DefaultUploaderConfig = &models.UploaderConfig{
	AllowedSources:         []string{models.FileUpload.String()},
	RequiredMetadataFields: []string{},
}

type WorkspaceService struct {
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
	wsUserRepo   *repo.WorkspaceUserRepo
	userRepo     *repo.UserRepo
}

func NewWorkspaceService(wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo,
	wsUserRepo *repo.WorkspaceUserRepo, userRepo *repo.UserRepo) *WorkspaceService {
	return &WorkspaceService{
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
		wsUserRepo:   wsUserRepo,
		userRepo:     userRepo,
	}
}

func (s *WorkspaceService) GetWorkspaces(ctx context.Context, userID string) ([]models.WorkspaceNameID, error) {
	workspaces, err := s.wsRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {

	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	workspace.CreatedBy = user.UserID
	workspace.UpdatedBy = user.UserID

	uw := &models.UserWorkspace{
		UserID:      user.UserID,
		WorkspaceID: workspace.ID,
		Role:        models.UserRoleOwner,
	}

	err = s.wsRepo.Create(ctx, workspace, uw, DefaultUploaderConfig)
	return err
}

func (s *WorkspaceService) GetWorkspaceDetails(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	ws, err := s.wsRepo.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]models.WorkspaceUserDetails, error) {
	users, err := s.wsUserRepo.GetUsersInWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *WorkspaceService) AddUserToWorkspace(ctx context.Context, workspaceID string, addReq *dto.AddWorkspaceUser) error {
	user, err := s.userRepo.GetByEmail(ctx, addReq.Email)
	if err != nil {
		infra.Log.Errorf("failed to get user by email: %s", err.Error())
		return fmt.Errorf(msg.UserNotFound, addReq.Email)
	}

	if addReq.Role != models.UserRoleContributor && addReq.Role != models.UserRoleViewer {
		return fmt.Errorf(msg.UnknownRole, addReq.Role)
	}

	uw := &models.UserWorkspace{
		UserID:      user.ID,
		WorkspaceID: workspaceID,
		Role:        addReq.Role,
	}
	err = s.wsUserRepo.AddUserToWorkspace(ctx, uw)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	isLastOwner, err := s.wsUserRepo.IsOwner(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if isLastOwner {
		return errors.New(msg.OwnerCannotBeRemoved)
	}
	err = s.wsUserRepo.RemoveUserFromWorkspace(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) ChangeUserRoleInWorkspace(ctx context.Context, workspaceID string, userID string, role models.UserRole) error {
	if role != models.UserRoleContributor && role != models.UserRoleViewer {
		return fmt.Errorf(msg.UnknownRole, role)
	}

	isOwner, err := s.wsUserRepo.IsOwner(ctx, workspaceID, userID)
	if err != nil {
		return err
	}

	if isOwner {
		return errors.New(msg.OwnerRoleCannotBeChanged)
	}

	uw := &models.UserWorkspace{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
	}

	err = s.wsUserRepo.UpdateUserInWorkspace(ctx, uw)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	config, err := s.wsConfigRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *WorkspaceService) SetUploaderConfig(ctx context.Context, workspaceID string, config *models.UploaderConfig) error {
	config.WorkspaceID = workspaceID
	err := s.wsConfigRepo.SetConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}
