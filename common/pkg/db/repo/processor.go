package repo

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type ProcessorRepo struct {
	db *db.DB
}

func NewProcessorRepo(db *db.DB) *ProcessorRepo {
	return &ProcessorRepo{
		db: db,
	}
}

func (r *ProcessorRepo) GetAll(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	var processors []models.Processor
	err := r.db.Orm.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("enabled desc, updated_at desc").
		Find(&processors).Error

	if err != nil {
		return nil, dbutils.DBError(err)
	}

	return processors, nil
}

func (r *ProcessorRepo) Get(ctx context.Context, processorID string) (*models.Processor, error) {
	var processor models.Processor
	err := r.db.Orm.WithContext(ctx).
		First(&processor, "id = ?", processorID).Error

	if err != nil {
		return nil, dbutils.DBError(err)
	}

	return &processor, nil
}

func (r *ProcessorRepo) Create(ctx context.Context, processor *models.Processor) error {
	err := r.db.Orm.WithContext(ctx).Create(processor).Error
	if err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *ProcessorRepo) Patch(ctx context.Context, workspaceID, processorID string, patch map[string]interface{}) error {
	if err := r.db.Orm.WithContext(ctx).Model(&models.Processor{}).Where("id = ?", processorID).Updates(patch).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *ProcessorRepo) Delete(ctx context.Context, workspaceID, processorID string) error {
	if err := r.db.Orm.WithContext(ctx).Delete(&models.Processor{}, "id = ?", processorID).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}

func (r *ProcessorRepo) SaveWorkflow(ctx context.Context, workspaceID, processorID string, workflow string) error {
	patch := map[string]interface{}{
		"workflow": workflow,
	}
	if err := r.db.Orm.WithContext(ctx).Model(&models.Processor{}).Where("id = ?", processorID).Updates(patch).Error; err != nil {
		return dbutils.DBError(err)
	}
	return nil
}
