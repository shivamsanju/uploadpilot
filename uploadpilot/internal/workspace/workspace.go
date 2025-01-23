package workspace

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/messages"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DefaultUploaderConfig = &models.UploaderConfig{
	AllowedSources:         []models.AllowedSources{models.FileUpload},
	RequiredMetadataFields: []string{},
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	workspace.UploaderConfig = DefaultUploaderConfig
	workspace.ID = primitive.NewObjectID()

	if err := s.validateWorkspaceBody(workspace); err != nil {
		return err
	}

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

func (s *WorkspaceService) GetUserWorkspaces(ctx context.Context, userID string) ([]models.Workspace, error) {
	workspaces, err := s.wsRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) validateWorkspaceBody(workspace *models.Workspace) error {
	if err := infra.Validate.Struct(workspace); err != nil {
		errs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errs[err.Field()] = err.Tag()
		}
		return fmt.Errorf(messages.ValidationErr, errs)
	}
	return nil
}
