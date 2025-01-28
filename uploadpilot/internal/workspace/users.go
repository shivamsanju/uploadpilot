package workspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
)

var ValidUserRoles = map[models.UserRole]bool{
	models.UserRoleOwner:       true,
	models.UserRoleContributor: true,
	models.UserRoleViewer:      true,
}

func (s *WorkspaceService) GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]models.WorkspaceUserWithDetails, error) {
	users, err := s.wsRepo.GetUsersInWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *WorkspaceService) AddUserToWorkspace(ctx context.Context, workspaceID string, addReq *dto.AddWorkspaceUser) error {
	if err := infra.Validator.ValidateBody(addReq); err != nil {
		return err
	}

	user, err := s.userRepo.GetByEmail(ctx, addReq.Email)
	if err != nil {
		infra.Log.Errorf("failed to get user by email: %s", err.Error())
		return fmt.Errorf(msg.UserNotFound, addReq.Email)
	}

	exists, err := s.wsRepo.CheckIfUserExistsInWorkspace(ctx, workspaceID, user.UserID)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf(msg.UserAlreadyExistsInWorkspace, user.Email)
	}

	if !ValidUserRoles[addReq.Role] {
		return fmt.Errorf(msg.UnknownRole, addReq.Role)
	}

	workspaceUser := &models.WorkspaceUser{
		UserID: user.UserID,
		Role:   addReq.Role,
	}

	err = s.wsRepo.AddUserToWorkspace(ctx, workspaceID, workspaceUser)
	return err
}

func (s *WorkspaceService) RemoveUserFromWorkspace(ctx context.Context, workspaceID string, userID string) error {
	isLastOwner, err := s.wsRepo.IsUserLastOwner(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if isLastOwner {
		return errors.New(msg.LastOwnerCannotBeRemoved)
	}
	return s.wsRepo.RemoveUserFromWorkspace(ctx, workspaceID, userID)
}

func (s *WorkspaceService) ChangeUserRoleInWorkspace(ctx context.Context, workspaceID string, userID string, role models.UserRole) error {
	if !ValidUserRoles[role] {
		return fmt.Errorf(msg.UnknownRole, role)
	}
	isLastOwner, err := s.wsRepo.IsUserLastOwner(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if isLastOwner && role == models.UserRoleOwner {
		return errors.New(msg.LastOwnerRoleCannotBeChanged)
	}
	return s.wsRepo.ChangeUserRoleInWorkspace(ctx, workspaceID, userID, role)
}
