package uploaderconfig

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type Service struct {
	configRepo *repo.WorkspaceConfigRepo
}

func NewConfigService(configRepo *repo.WorkspaceConfigRepo) *Service {
	return &Service{
		configRepo: configRepo,
	}
}
func (s *Service) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	config, err := s.configRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}
