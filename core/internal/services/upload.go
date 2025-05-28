package services

import (
	"context"
	"fmt"
	"slices"
	"time"

	"maps"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/msg"
	"github.com/uploadpilot/core/internal/rbac"
	"github.com/uploadpilot/core/pkg/utils"
	"github.com/uploadpilot/core/web/webutils"
)

type UploadService struct {
	accessManager *rbac.AccessManager
	uploadRepo    *repo.UploadRepo
	workspaceSvc  *WorkspaceService
	processorSvc  *ProcessorService
	s3Client      *s3.Client
}

func NewUploadService(accessManager *rbac.AccessManager, uploadRepo *repo.UploadRepo, workspaceSvc *WorkspaceService, processorSvc *ProcessorService, s3Client *s3.Client) *UploadService {
	return &UploadService{
		accessManager,
		uploadRepo,
		workspaceSvc,
		processorSvc,
		s3Client,
	}
}

// CreateUpload creates a new upload
func (s *UploadService) CreateUpload(ctx context.Context, tenantID, workspaceID string, upload *dto.CreateUploadRequest) (*dto.CreateUploadResponse, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Uploader) {
		return nil, fmt.Errorf(msg.ErrAccessDenied)
	}

	if err := s.workspaceSvc.ValidateUploadRequestWithConfig(ctx, tenantID, workspaceID, upload); err != nil {
		return nil, err
	}
	newUploadID := uuid.New().String()
	s3CompatibleFileName := utils.ConvertToS3CompatibleFilename(upload.FileName)
	objectKey := fmt.Sprintf("%s/raw/%s", newUploadID, s3CompatibleFileName)
	req, err := s.createSingleUseSignedUploadURL(workspaceID, objectKey, upload)
	if err != nil {
		return nil, err
	}

	if err := s.uploadRepo.Create(ctx, &models.Upload{
		ID:            newUploadID,
		WorkspaceID:   workspaceID,
		FileName:      s3CompatibleFileName,
		ContentType:   upload.ContentType,
		ContentLength: upload.ContentLength,
		Metadata:      upload.Metadata,
		StartedAt:     time.Now(),
		Status:        models.UploasStatusInProgress,
	}); err != nil {
		return nil, err
	}

	signedHeaders := make(map[string][]string)
	maps.Copy(signedHeaders, req.SignedHeader)
	return &dto.CreateUploadResponse{
		UploadID:      newUploadID,
		UploadURL:     req.URL,
		Method:        req.Method,
		SignedHeaders: signedHeaders,
	}, nil
}

func (s *UploadService) createSingleUseSignedUploadURL(bucketName, objectKey string, uploadReq *dto.CreateUploadRequest) (*v4.PresignedHTTPRequest, error) {
	log.Debug().Interface("uploadReq", uploadReq).Str("bucketName", bucketName).Str("objectKey", objectKey).Msg("upload request")
	request, err := s3.NewPresignClient(s.s3Client).PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(uploadReq.ContentType),
		ContentLength: aws.Int64(uploadReq.ContentLength),
		IfNoneMatch:   aws.String("*"),
	}, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(uploadReq.UploadURLValiditySecs) * time.Second
	})
	if err != nil {
		errID := uuid.New().String()
		log.Error().Err(err).Str("errID", errID).Msg("failed to create single use signed upload url")
		return nil, fmt.Errorf(msg.ErrUnexpected, errID)
	}
	return request, err
}

// FinishUpload finishes an upload
func (s *UploadService) FinishUpload(ctx context.Context, tenantID, workspaceID, uploadID string, status models.UploadStatus) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Uploader) {
		return fmt.Errorf(msg.ErrAccessDenied)
	}

	upload, err := s.uploadRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	if !slices.Contains(models.UploadNonTerminalStates, upload.Status) {
		return fmt.Errorf(msg.ErrUploadAlreadyIsTerminalState)
	}

	upload.FinishedAt = time.Now()
	upload.Status = status

	if err = s.processorSvc.TriggerWorkflows(ctx, workspaceID, upload); err != nil {
		log.Error().Str("workspace_id", workspaceID).Str("upload_id", uploadID).Err(err).Msg("failed to trigger workflows")
		upload.Status = models.UploadStatusFailed
		s.uploadRepo.Update(ctx, uploadID, upload)
		return err
	}

	if err := s.uploadRepo.Update(ctx, uploadID, upload); err != nil {
		return err
	}

	return nil
}

// GetAllUploadsForWorkspace returns all uploads
func (s *UploadService) GetAllUploadsForWorkspace(ctx context.Context, tenantID, workspaceID string, paginationParams *models.PaginationParams) ([]models.Upload, int64, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, 0, err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return nil, 0, fmt.Errorf(msg.ErrAccessDenied)
	}
	return s.uploadRepo.GetAll(ctx, workspaceID, paginationParams)
}

func (s *UploadService) GetUploadDetails(ctx context.Context, tenantID, workspaceID, uploadID string) (*models.Upload, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return nil, fmt.Errorf(msg.ErrAccessDenied)
	}

	return s.uploadRepo.Get(ctx, uploadID)
}

func (s *UploadService) ProcessUpload(ctx context.Context, tenantID, workspaceID, uploadID string) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return fmt.Errorf(msg.ErrAccessDenied)
	}

	upload, err := s.uploadRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	return s.processorSvc.TriggerWorkflows(ctx, workspaceID, upload)
}

func (s *UploadService) DeleteUpload(ctx context.Context, tenantID, workspaceID, uploadID string) error {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Admin) {
		return fmt.Errorf(msg.ErrAccessDenied)
	}

	return s.uploadRepo.Delete(ctx, uploadID)
}

func (s *UploadService) GetUploadSignedURL(ctx context.Context, tenantID, workspaceID, uploadID string) (string, error) {
	session, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return "", err
	}

	if !s.accessManager.CheckAccess(session.Sub, tenantID, workspaceID, rbac.Reader) {
		return "", fmt.Errorf(msg.ErrAccessDenied)
	}

	upload, err := s.uploadRepo.Get(ctx, uploadID)
	if err != nil {
		return "", err
	}
	if upload.Status != models.UploadStatusFinished {
		return "", fmt.Errorf(msg.ErrUploadNotFinished)
	}

	objectKey := fmt.Sprintf("%s/raw/%s", uploadID, upload.FileName)

	expiry := time.Now().Add(15 * time.Minute)
	resp, err := s3.NewPresignClient(s.s3Client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:          &workspaceID,
		Key:             &objectKey,
		ResponseExpires: &expiry,
	})
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
