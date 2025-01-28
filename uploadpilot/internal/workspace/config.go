package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

var DefaultUploaderConfig = &models.UploaderConfig{
	AllowedSources:         []models.AllowedSources{models.FileUpload},
	RequiredMetadataFields: []string{},
}

func (s *WorkspaceService) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	config, err := s.wsRepo.GetUploaderConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *WorkspaceService) SetUploaderConfig(ctx context.Context, workspaceID string, config *models.UploaderConfig) error {
	if err := infra.Validator.ValidateBody(config); err != nil {
		return err
	}
	err := s.wsRepo.SetUploaderConfig(ctx, workspaceID, config)
	if err != nil {
		return err
	}
	return nil
}
