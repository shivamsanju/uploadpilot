package svc

import (
	"context"
	"strings"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type UploadService struct {
	upRepo       *repo.UploadRepo
	wsRepo       *repo.WorkspaceRepo
	wsConfigRepo *repo.WorkspaceConfigRepo
	userRepo     *repo.UserRepo
	logRepo      *repo.UploadLogsRepo
}

func NewUploadService(upRepo *repo.UploadRepo, wsRepo *repo.WorkspaceRepo, wsConfigRepo *repo.WorkspaceConfigRepo,
	userRepo *repo.UserRepo, logRepo *repo.UploadLogsRepo) *UploadService {
	return &UploadService{
		upRepo:       upRepo,
		wsRepo:       wsRepo,
		wsConfigRepo: wsConfigRepo,
		userRepo:     userRepo,
		logRepo:      logRepo,
	}
}

func (us *UploadService) GetAllUploads(ctx context.Context, workspaceID string, skip int, limit int, search string) ([]models.Upload, int64, error) {
	if strings.HasPrefix(search, "{") {
		searchParams, err := utils.ExtractKeyValuePairs(search)
		if err != nil {
			return nil, 0, err
		}
		return us.upRepo.GetAllFilterByMetadata(ctx, workspaceID, skip, limit, searchParams)
	}

	return us.upRepo.GetAll(ctx, workspaceID, skip, limit, search)
}

func (us *UploadService) GetUploadDetails(ctx context.Context, workspaceID, uploadID string) (*models.Upload, error) {
	return us.upRepo.Get(ctx, uploadID)
}

func (us *UploadService) GetLogs(ctx context.Context, uploadID string) ([]models.UploadLog, error) {
	logs, err := us.logRepo.GetLogs(ctx, uploadID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
