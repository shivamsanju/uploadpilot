package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type ProcessorRepo struct {
}

func NewProcessorRepo() *ProcessorRepo {
	return &ProcessorRepo{}
}

func (i *ProcessorRepo) GetAll(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	var processors []models.Processor

	if err := sqlDB.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("updated_at desc").Find(&processors).Error; err != nil {
		return nil, utils.DBError(err)
	}

	return processors, nil
}

func (i *ProcessorRepo) Get(ctx context.Context, processorID string) (*models.Processor, error) {
	var processor models.Processor

	if err := sqlDB.WithContext(ctx).First(&processor, "id = ?", processorID).Error; err != nil {
		return nil, utils.DBError(err)
	}

	return &processor, nil
}

func (i *ProcessorRepo) Create(ctx context.Context, processor *models.Processor) error {
	if err := sqlDB.WithContext(ctx).Create(processor).Error; err != nil {
		infra.Log.Errorf("failed to create processor: %s", err.Error())
		return utils.DBError(err)
	}

	return nil
}

func (i *ProcessorRepo) Patch(ctx context.Context, processorID string, patch map[string]interface{}) error {
	if err := sqlDB.WithContext(ctx).Model(&models.Processor{}).Where("id = ?", processorID).Updates(patch).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}

func (i *ProcessorRepo) Delete(ctx context.Context, processorID string) error {
	if err := sqlDB.WithContext(ctx).Delete(&models.Processor{}, "id = ?", processorID).Error; err != nil {
		return utils.DBError(err)
	}

	return nil
}
