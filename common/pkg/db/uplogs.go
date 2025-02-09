package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/dbutils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type UploadLogsRepo struct {
}

func NewUploadLogsRepo() *UploadLogsRepo {
	return &UploadLogsRepo{}
}

func (u *UploadLogsRepo) GetLogs(ctx context.Context, uploadID string) ([]models.UploadLog, error) {
	var logs []models.UploadLog
	if err := sqlDB.WithContext(ctx).Where("upload_id = ?", uploadID).Find(&logs).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return logs, nil
}

func (u *UploadLogsRepo) BatchAddLogs(ctx context.Context, logs []*models.UploadLog) error {
	err := sqlDB.WithContext(ctx).Create(logs).Error
	if err != nil {
		return dbutils.DBError(err)
	}
	return nil
}
