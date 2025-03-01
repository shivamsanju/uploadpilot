package upload

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/phuslu/log"
	"github.com/uploadpilot/manager/internal/config"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/db/repo"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/manager/internal/svc/processor"
)

type Service struct {
	upRepo       *repo.UploadRepo
	processorSvc *processor.Service
}

func NewService(upRepo *repo.UploadRepo, processorSvc *processor.Service) *Service {
	return &Service{
		upRepo:       upRepo,
		processorSvc: processorSvc,
	}
}

func (us *Service) GetAllUploads(ctx context.Context, workspaceID string, paginationParams *models.PaginationParams) ([]models.Upload, int64, error) {
	return us.upRepo.GetAll(ctx, workspaceID, paginationParams)
}

func (us *Service) GetUploadDetails(ctx context.Context, workspaceID, uploadID string) (*models.Upload, error) {
	return us.upRepo.Get(ctx, uploadID)
}

// Upload Related Methods
func (us *Service) CreateUpload(ctx context.Context, workspaceID string, upload *models.Upload) error {
	upload.Status = models.UploadStatusInProgress
	upload.WorkspaceID = workspaceID
	return us.upRepo.Create(ctx, workspaceID, upload)
}

func (us *Service) FinishUpload(ctx context.Context, workspaceID, uploadID string, req *dto.FinishUploadRequest) error {
	upload, err := us.upRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	upload.FinishedAt = req.FinishedAt
	if req.Status != "" {
		upload.Status = models.UploadStatus(req.Status)
	}
	if req.Size != 0 {
		upload.Size = req.Size
	}
	err = us.processorSvc.TriggerWorkflows(ctx, workspaceID, upload)
	if err != nil {
		log.Error().Err(err).Msg("failed to trigger workflows")
		upload.Status = models.UploadStatusFailed
		upload.StatusReason = "failed to trigger workflows"
		us.upRepo.Update(ctx, uploadID, upload)
		return err
	}
	return us.upRepo.Update(ctx, uploadID, upload)
}

func (us *Service) ProcessUpload(ctx context.Context, workspaceID, uploadID string) error {
	upload, err := us.upRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	return us.processorSvc.TriggerWorkflows(ctx, workspaceID, upload)
}
func (us *Service) DeleteUpload(ctx context.Context, workspaceID, uploadID string) error {
	return us.upRepo.Delete(ctx, uploadID)
}

func (us *Service) GetUploadSignedURL(ctx context.Context, workspaceID, uploadID string) (string, error) {
	expiry := time.Now().Add(15 * time.Minute)
	resp, err := s3.NewPresignClient(infra.S3Client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:          &config.S3BucketName,
		Key:             &uploadID,
		ResponseExpires: &expiry,
	})
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
