package repo

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/uploadpilot/core/internal/db/driver"
	"github.com/uploadpilot/core/internal/db/models"
	dbutils "github.com/uploadpilot/core/internal/db/utils"
)

type UploadRepo struct {
	db *driver.Driver
}

func NewUploadRepo(db *driver.Driver) *UploadRepo {
	return &UploadRepo{
		db: db,
	}
}

func (r *UploadRepo) GetAll(ctx context.Context, workspaceID string, paginationParams *models.PaginationParams) ([]models.Upload, int64, error) {
	var uploads []models.Upload

	query := r.db.Orm.WithContext(ctx).
		Model(&models.Upload{}).
		Select("id", "file_name", "status", "file_type", "started_at", "size", "finished_at").
		Where("workspace_id = ?", workspaceID)

	query, totalRecords, sortApplied, err := dbutils.BuildPaginationQuery(
		query,
		&dbutils.PaginationQueryInput{
			PaginationParams:    paginationParams,
			AllowedSearchFields: []string{"file_name", "status", "file_type"},
			AllowedFilterFields: []string{"status"},
		},
	)

	if err != nil {
		return nil, 0, err
	}

	if !sortApplied {
		query = query.Order("finished_at DESC")
	}

	if err := query.Find(&uploads).Error; err != nil {
		return nil, 0, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return uploads, totalRecords, nil
}

func (r *UploadRepo) Get(ctx context.Context, uploadID string) (*models.Upload, error) {
	var upload models.Upload
	if err := r.db.Orm.WithContext(ctx).First(&upload, "id = ?", uploadID).Error; err != nil {
		return nil, dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return &upload, nil
}

func (r *UploadRepo) Create(ctx context.Context, workspaceID string, upload *models.Upload) error {
	if err := r.db.Orm.WithContext(ctx).Create(upload).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}
	return nil
}

func (r *UploadRepo) Update(ctx context.Context, uploadID string, upload *models.Upload) error {
	if err := r.db.Orm.WithContext(ctx).Save(upload).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return nil
}

func (r *UploadRepo) Delete(ctx context.Context, uploadID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Upload{}, "id = ?", uploadID).Error; err != nil {
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
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
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
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
		return dbutils.DBError(ctx, r.db.Orm.Logger, err)
	}

	return nil
}
