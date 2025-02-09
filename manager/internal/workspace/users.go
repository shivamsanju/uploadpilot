package workspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/msg"
)

func (s *WorkspaceService) GetWorkspaceUsers(ctx context.Context, workspaceID string) ([]dto.WorkspaceUser, error) {
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
