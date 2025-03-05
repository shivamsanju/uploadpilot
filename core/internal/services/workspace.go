package services

import (
	"context"
	"fmt"

	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/web/webutils"
)

type WorkspaceService struct {
	acm          *rbac.AccessManager
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
}

func NewWorkspaceService(accessManager *rbac.AccessManager, wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo) *WorkspaceService {
	return &WorkspaceService{
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
		acm:          accessManager,
	}
}

func (s *WorkspaceService) GetAllWorkspaces(ctx context.Context) ([]models.Workspace, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	access, err := s.acm.GetSubjectTenantAccess(session.UserID, session.TenantID)
	if err != nil {
		return nil, err
	}

	log.Debug().Interface("access", access).Msg("access")

	workspaces, err := s.wsRepo.GetAll(ctx, session.TenantID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspaceInfo(ctx context.Context, workspaceID string) (*models.Workspace, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.acm.CheckAccess(session.UserID, session.TenantID, workspaceID, rbac.Reader) {
		return nil, err
	}

	ws, err := s.wsRepo.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, workspace *models.Workspace) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.acm.CheckAccess(session.UserID, session.TenantID, "*", rbac.Admin) {
		return err
	}

	if err := s.acm.AddAccess(session.UserID, session.TenantID, workspace.ID, rbac.Admin); err != nil {
		return err
	}

	workspace.CreatedBy = session.UserID
	workspace.UpdatedBy = session.UserID
	workspace.TenantID = session.TenantID
	err = s.wsRepo.Create(ctx, workspace)
	return err
}

func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, workspaceID string) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.acm.CheckAccess(session.UserID, session.TenantID, workspaceID, rbac.Admin) {
		return err
	}

	return s.wsRepo.Delete(ctx, workspaceID)
}

func (s *WorkspaceService) GetWorkspaceConfig(ctx context.Context, workspaceID string) (*models.WorkspaceConfig, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	hascess := s.acm.CheckAccess(session.UserID, session.TenantID, workspaceID, rbac.Reader)
	log.Debug().Str("UserID", session.UserID).Str("TenantID", session.TenantID).Str("WorkspaceID", workspaceID).Bool("hascess", hascess).Msg("GetWorkspaceConfig")
	if !hascess {
		return nil, fmt.Errorf("access denied")
	}
	config, err := s.wsConfigRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *WorkspaceService) SetWorkspaceConfig(ctx context.Context, workspaceID string, config *models.WorkspaceConfig) error {
	config.WorkspaceID = workspaceID
	err := s.wsConfigRepo.SetConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}
