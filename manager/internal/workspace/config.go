package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
)

var DefaultUploaderConfig = &models.UploaderConfig{
	AllowedSources:         []string{models.FileUpload.String()},
	RequiredMetadataFields: []string{},
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
