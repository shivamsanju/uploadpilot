package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/web/webutils"
)

type WorkspaceService struct {
	acm          *rbac.AccessManager
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
	s3Client     *s3.Client
}

func NewWorkspaceService(accessManager *rbac.AccessManager, wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo, s3Client *s3.Client) *WorkspaceService {
	return &WorkspaceService{
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
		acm:          accessManager,
		s3Client:     s3Client,
	}
}

func (s *WorkspaceService) GetAllWorkspaces(ctx context.Context, tenantID string) ([]models.Workspace, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	access, err := s.acm.GetSubjectTenantAccess(session.UserID, tenantID)
	if err != nil {
		return nil, err
	}

	log.Debug().Interface("access", access).Msg("access")

	workspaces, err := s.wsRepo.GetAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspaceInfo(ctx context.Context, tenantID, workspaceID string) (*models.Workspace, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.acm.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return nil, err
	}

	ws, err := s.wsRepo.Get(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, tenantID string, workspace *models.Workspace) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.acm.CheckAccess(session.Sub, tenantID, "*", rbac.Admin) {
		return err
	}

	workspace.ID = uuid.New().String()
	workspace.CreatedBy = session.UserID
	workspace.UpdatedBy = session.UserID
	workspace.TenantID = tenantID

	if err := s.acm.AddAccess(session.UserID, tenantID, workspace.ID, rbac.Admin); err != nil {
		return err
	}

	_, err = s.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket:                    &workspace.ID,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraint(s.s3Client.Options().Region)},
	})
	if err != nil {
		errID := uuid.New().String()
		log.Error().Err(err).Str("errID", errID).Msg("failed to create bucket for workspace.")
		return fmt.Errorf(msg.ErrUnexpected, errID)
	}

	err = s.wsRepo.Create(ctx, workspace)
	return err
}

func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, tenantID, workspaceID string) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.acm.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Admin) {
		return err
	}

	return s.wsRepo.Delete(ctx, workspaceID)
}

func (s *WorkspaceService) GetWorkspaceConfig(ctx context.Context, tenantID, workspaceID string) (*models.WorkspaceConfig, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	if !s.acm.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return nil, fmt.Errorf(msg.ErrAccessDenied)
	}
	config, err := s.wsConfigRepo.GetConfig(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *WorkspaceService) SetWorkspaceConfig(ctx context.Context, tenantID, workspaceID string, config *models.WorkspaceConfig) error {
	config.WorkspaceID = workspaceID
	err := s.wsConfigRepo.SetConfig(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (s *WorkspaceService) ValidateUploadRequestWithConfig(ctx context.Context, tenantID string, workspaceID string, uploadReq *dto.CreateUploadRequest) error {
	config, err := s.GetWorkspaceConfig(ctx, tenantID, workspaceID)
	if err != nil {
		return err
	}

	if config.MinFileSize != nil && *config.MinFileSize != 0 && uploadReq.ContentLength < *config.MinFileSize {
		return fmt.Errorf(msg.ErrUploadSizeLessThanRequired, uploadReq.ContentLength, *config.MinFileSize)
	}

	if config.MaxFileSize != nil && *config.MaxFileSize != 0 && uploadReq.ContentLength > *config.MaxFileSize {
		return fmt.Errorf(msg.ErrUploadSizeExceedAllowedLimit, uploadReq.ContentLength, *config.MaxFileSize)
	}

	if len(config.AllowedContentTypes) > 0 && !slices.Contains(config.AllowedContentTypes, uploadReq.ContentType) {
		return fmt.Errorf(msg.ErrUploadContentTypeNotAllowed)
	}

	if config.MaxUploadURLLifetimeSecs != 0 && uploadReq.UploadURLValiditySecs > config.MaxUploadURLLifetimeSecs {
		return fmt.Errorf(msg.ErrUploadURLValidityExceedsAllowedLimit, uploadReq.UploadURLValiditySecs, config.MaxUploadURLLifetimeSecs)
	}

	if len(config.RequiredMetadataFields) > 0 {
		missingFields := ""
		for _, field := range config.RequiredMetadataFields {
			if _, ok := uploadReq.Metadata[field]; !ok {
				missingFields = fmt.Sprintf("%s, %s", missingFields, field)
			}
		}
		if len(missingFields) > 0 {
			return fmt.Errorf(msg.ErrUploadMissingRequiredMetadataFields, missingFields)
		}
	}

	return nil
}

func (s *WorkspaceService) GetTenantID(ctx context.Context, workspaceID string) (string, error) {
	return s.wsRepo.GetTenantID(ctx, workspaceID)
}
