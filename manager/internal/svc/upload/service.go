package upload

import (
	"context"
	"strings"

	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/utils"
)

type Service struct {
	upRepo       *repo.UploadRepo
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
	userRepo     *repo.UserRepo
	logRepo      *repo.UploadLogsRepo
}

func NewService(upRepo *repo.UploadRepo, wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo,
	userRepo *repo.UserRepo, logRepo *repo.UploadLogsRepo) *Service {
	return &Service{
		upRepo:       upRepo,
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
		userRepo:     userRepo,
		logRepo:      logRepo,
	}
}

func (us *Service) GetAllUploads(ctx context.Context, workspaceID string, skip int, limit int, search string) ([]models.Upload, int64, error) {
	if strings.HasPrefix(search, "{") {
		searchParams, err := utils.ExtractKeyValuePairs(search)
		if err != nil {
			return nil, 0, err
		}
		return us.upRepo.GetAllFilterByMetadata(ctx, workspaceID, skip, limit, searchParams)
	}

	return us.upRepo.GetAll(ctx, workspaceID, skip, limit, search)
}

func (us *Service) GetUploadDetails(ctx context.Context, workspaceID, uploadID string) (*models.Upload, error) {
	return us.upRepo.Get(ctx, uploadID)
}

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
	return us.upRepo.Update(ctx, uploadID, upload)
}
