package workspace

import (
	"context"

	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/db/repo"
	"github.com/uploadpilot/manager/internal/utils"
)

type Service struct {
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
}

func NewService(wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo) *Service {
	return &Service{
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
	}
}

func (s *Service) GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	workspaces, err := s.wsRepo.GetAll(ctx, session.TenantID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *Service) GetWorkspaceInfo(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	ws, err := s.wsRepo.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *Service) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	workspace.CreatedBy = session.UserID
	workspace.UpdatedBy = session.UserID
	workspace.TenantID = session.TenantID
	err = s.wsRepo.Create(ctx, workspace)
	return err
}

func (s *Service) DeleteWorkspace(ctx context.Context, workspaceID string) error {
	return s.wsRepo.Delete(ctx, workspaceID)
}

func (s *Service) GetWorkspaceConfig(ctx context.Context, workspaceID string) (*models.WorkspaceConfig, error) {
	config, err := s.wsConfigRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *Service) SetWorkspaceConfig(ctx context.Context, workspaceID string, config *models.WorkspaceConfig) error {
	config.WorkspaceID = workspaceID
	err := s.wsConfigRepo.SetConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}
