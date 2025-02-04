package db

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type UploadRepo struct {
}

func NewUploadRepo() *UploadRepo {
	return &UploadRepo{}
}

func (i *UploadRepo) GetAll(ctx context.Context, workspaceID string, skip, limit int, search string) ([]models.Upload, int64, error) {
	var uploads []models.Upload
	var totalRecords int64

	query := sqlDB.WithContext(ctx).Model(&models.Upload{}).Where("workspace_id = ?", workspaceID)

	if search != "" {
		query = query.Where("name LIKE ? OR status LIKE ? OR stored_file_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("finished_at DESC").Offset(skip).Limit(limit).Find(&uploads).Error; err != nil {
		return nil, 0, err
	}

	return uploads, totalRecords, nil
}

func (i *UploadRepo) GetAllFilterByMetadata(ctx context.Context, workspaceID string, skip, limit int, search map[string]string) ([]models.Upload, int64, error) {
	var uploads []models.Upload
	var totalRecords int64

	query := sqlDB.WithContext(ctx).Model(&models.Upload{}).Where("workspace_id = ?", workspaceID)

	for key, value := range search {
		if key != "" && value != "" {
			query = query.Where("metadata->>? = ?", key, value)
		}
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("finished_at DESC").Offset(int(skip)).Limit(int(limit)).Find(&uploads).Error; err != nil {
		return nil, 0, err
	}

	return uploads, totalRecords, nil
}

func (i *UploadRepo) Get(ctx context.Context, uploadID string) (*models.Upload, error) {
	var upload models.Upload
	if err := sqlDB.WithContext(ctx).First(&upload, "id = ?", upload).Error; err != nil {
		return nil, utils.DBError(err)
	}
	return &upload, nil
}

func (i *UploadRepo) Create(ctx context.Context, workspaceID string, upload *models.Upload) error {
	if err := sqlDB.WithContext(ctx).Create(upload).Error; err != nil {
		return utils.DBError(err)
	}
	return nil
}

func (i *UploadRepo) Update(ctx context.Context, uploadID string, upload *models.Upload) error {
	if err := sqlDB.WithContext(ctx).Save(upload).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}

func (i *UploadRepo) Delete(ctx context.Context, uploadID string) error {
	if err := sqlDB.WithContext(ctx).Delete(&models.Upload{}, "id = ?", uploadID).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}

func (i *UploadRepo) SetStatus(ctx context.Context, uploadID string, status models.UploadStatus) error {
	update := map[string]interface{}{
		"status": status,
	}

	if slices.Contains(models.UploadTerminalStates, status) {
		update["finished_at"] = time.Now()
	}

	if err := sqlDB.WithContext(ctx).Model(&models.Upload{}).Where("id = ?", uploadID).Updates(update).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}

func (i *UploadRepo) Patch(ctx context.Context, uploadID string, patchMap map[string]interface{}) error {
	patch := map[string]interface{}{}
	for key, value := range patchMap {
		if value == nil {
			delete(patchMap, key)
		}

		if !slices.Contains([]string{"url", "processedUrl", "storedFileName"}, key) {
			return fmt.Errorf("unsupported patch key: %s", key)
		}

		patch[key] = value
	}

	if err := sqlDB.WithContext(ctx).Model(&models.Upload{}).Where("id = ?", uploadID).Updates(patch).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}
