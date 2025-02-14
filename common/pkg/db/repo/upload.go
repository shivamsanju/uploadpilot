package repo

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type UploadRepo struct {
	db *db.DB
}

func NewUploadRepo(db *db.DB) *UploadRepo {
	return &UploadRepo{
		db: db,
	}
}

func (r *UploadRepo) GetAll(ctx context.Context, workspaceID string, skip, limit int, search string) ([]models.Upload, int64, error) {
	var uploads []models.Upload
	var totalRecords int64

	query := r.db.Orm.WithContext(ctx).Model(&models.Upload{}).Where("workspace_id = ?", workspaceID)

	if search != "" {
		query = query.Where("name LIKE ? OR status LIKE ? OR stored_file_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, dbutils.DBError(err)
	}

	if err := query.Order("finished_at DESC").Offset(skip).Limit(limit).Find(&uploads).Error; err != nil {
		return nil, 0, dbutils.DBError(err)
	}

	return uploads, totalRecords, nil
}

func (r *UploadRepo) GetAllFilterByMetadata(ctx context.Context, workspaceID string, skip, limit int, search map[string]string) ([]models.Upload, int64, error) {
	var uploads []models.Upload
	var totalRecords int64

	query := r.db.Orm.WithContext(ctx).Model(&models.Upload{}).Where("workspace_id = ?", workspaceID)

	for key, value := range search {
		if key != "" && value != "" {
			query = query.Where("metadata->>? = ?", key, value)
		}
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, dbutils.DBError(err)
	}

	if err := query.Order("finished_at DESC").Offset(int(skip)).Limit(int(limit)).Find(&uploads).Error; err != nil {
		return nil, 0, dbutils.DBError(err)
	}

	return uploads, totalRecords, nil
}

func (r *UploadRepo) Get(ctx context.Context, uploadID string) (*models.Upload, error) {
	var upload models.Upload
	if err := r.db.Orm.WithContext(ctx).First(&upload, "id = ?", uploadID).Error; err != nil {
		return nil, dbutils.DBError(err)
	}
	return &upload, nil
}

func (r *UploadRepo) Create(ctx context.Context, workspaceID string, upload *models.Upload) error {
	if err := r.db.Orm.WithContext(ctx).Create(upload).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *UploadRepo) Update(ctx context.Context, uploadID string, upload *models.Upload) error {
	if err := r.db.Orm.WithContext(ctx).Save(upload).Error; err != nil {
		return dbutils.DBError(err)
	}

	return nil
}

func (r *UploadRepo) Delete(ctx context.Context, uploadID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Upload{}, "id = ?", uploadID).Error; err != nil {
		return dbutils.DBError(err)
	}

	return nil
}

func (r *UploadRepo) SetStatus(ctx context.Context, uploadID string, status models.UploadStatus) error {
	update := map[string]interface{}{
		"status": status,
	}

	if slices.Contains(models.UploadTerminalStates, status) {
		update["finished_at"] = time.Now()
	}

	if err := r.db.Orm.WithContext(ctx).Model(&models.Upload{}).Where("id = ?", uploadID).Updates(update).Error; err != nil {
		return dbutils.DBError(err)
	}

	return nil
}

func (r *UploadRepo) Patch(ctx context.Context, uploadID string, patchMap map[string]interface{}) error {
	patch := map[string]interface{}{}
	for key, value := range patchMap {
		if value == nil {
			delete(patchMap, key)
		}

		if !slices.Contains([]string{"url", "processed_url", "stored_file_name"}, key) {
			return fmt.Errorf("unsupported patch key: %s", key)
		}

		patch[key] = value
	}

	if err := r.db.Orm.WithContext(ctx).Model(&models.Upload{}).Where("id = ?", uploadID).Updates(patch).Error; err != nil {
		return dbutils.DBError(err)
	}

	return nil
}
