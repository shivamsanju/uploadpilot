package svc

import (
	"context"
	"strings"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type UploadService struct {
	upRepo       *db.UploadRepo
	wsRepo       *db.WorkspaceRepo
	wsConfigRepo *db.WorkspaceConfigRepo
	userRepo     *db.UserRepo
	logRepo      *db.UploadLogsRepo
}

func NewUploadService() *UploadService {
	return &UploadService{
		upRepo:       db.NewUploadRepo(),
		wsRepo:       db.NewWorkspaceRepo(),
		wsConfigRepo: db.NewWorkspaceConfigRepo(),
		userRepo:     db.NewUserRepo(),
		logRepo:      db.NewUploadLogsRepo(),
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
