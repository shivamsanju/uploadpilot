package workspace

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
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
	if err := infra.Validate.Struct(config); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return fmt.Errorf(msg.ValidationErr, errors)
	}
	err := s.wsRepo.SetUploaderConfig(ctx, workspaceID, config)
	if err != nil {
		return err
	}
	return nil
}
