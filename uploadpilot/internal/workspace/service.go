package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkspaceService struct {
	wsRepo   *db.WorkspaceRepo
	userRepo *db.UserRepo
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		wsRepo:   db.NewWorkspaceRepo(),
		userRepo: db.NewUserRepo(),
	}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	workspace.UploaderConfig = DefaultUploaderConfig
	workspace.ID = primitive.NewObjectID()

	_, err := s.wsRepo.Create(ctx, workspace)
	return err
}

func (s *WorkspaceService) GetWorkspaceDetails(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	ws, err := s.wsRepo.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) GetUserWorkspaces(ctx context.Context, userID string) ([]dto.WorkspaceNameID, error) {
	workspaces, err := s.wsRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}
