package uploaderconfig

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type UploaderConfigService struct {
	configRepo *db.WorkspaceConfigRepo
}

func NewUploaderConfigService() *UploaderConfigService {
	return &UploaderConfigService{
		configRepo: db.NewWorkspaceConfigRepo(),
	}
}
func (s *UploaderConfigService) GetUploaderConfig(ctx context.Context, workspaceID string) (*models.UploaderConfig, error) {
	config, err := s.configRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}
