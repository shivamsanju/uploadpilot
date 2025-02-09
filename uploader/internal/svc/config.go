package svc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

func (s *Service) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	config, err := s.configRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}
