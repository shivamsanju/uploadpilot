package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type UploadLogsRepo struct {
	db *db.DB
}

func NewUploadLogsRepo(db *db.DB) *UploadLogsRepo {
	return &UploadLogsRepo{
		db: db,
	}
}

func (r *UploadLogsRepo) GetLogs(ctx context.Context, uploadID string) ([]models.UploadLog, error) {
	var logs []models.UploadLog
	if err := r.db.Orm.WithContext(ctx).Where("upload_id = ?", uploadID).Find(&logs).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return logs, nil
}

func (r *UploadLogsRepo) BatchAddLogs(ctx context.Context, logs []*models.UploadLog) error {
	err := r.db.Orm.WithContext(ctx).Create(logs).Error
	if err != nil {
		return dbutils.DBError(err)
	}
	return nil
}
