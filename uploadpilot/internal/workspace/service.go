package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type WorkspaceService struct {
	wsRepo       *db.WorkspaceRepo
	wsConfigRepo *db.WorkspaceConfigRepo
	wsUserRepo   *db.WorkspaceUserRepo
	userRepo     *db.UserRepo
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		wsRepo:       db.NewWorkspaceRepo(),
		wsConfigRepo: db.NewWorkspaceConfigRepo(),
		wsUserRepo:   db.NewWorkspaceUserRepo(),
		userRepo:     db.NewUserRepo(),
	}
}

func (s *WorkspaceService) GetWorkspaces(ctx context.Context, userID string) ([]dto.WorkspaceNameID, error) {
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
