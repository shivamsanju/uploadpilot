package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/cache"
	"github.com/uploadpilot/uploadpilot/common/pkg/db/dbutils"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type ProcessorRepo struct {
}

func NewProcessorRepo() *ProcessorRepo {
	return &ProcessorRepo{}
}

func (i *ProcessorRepo) GetAll(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	dbFetchFn := func(processors *[]models.Processor) error {
		return sqlDB.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("enabled desc, updated_at desc").Find(processors).Error
	}

	var processors []models.Processor
	cl := cache.NewClient[*[]models.Processor](0)
	key := WorkspaceProcessorsKey(workspaceID)

	if err := cl.Query(ctx, key, &processors, dbFetchFn); err != nil {
		return nil, err
	}

	return processors, nil
}

func (i *ProcessorRepo) Get(ctx context.Context, processorID string) (*models.Processor, error) {
	dbFetchFn := func(processor *models.Processor) error {
		return sqlDB.WithContext(ctx).
			First(processor, "id = ?", processorID).Error
	}

	var processor models.Processor
	cl := cache.NewClient[*models.Processor](0)
	key := ProcessorKey(processorID)

	if err := cl.Query(ctx, key, &processor, dbFetchFn); err != nil {
		return nil, err
	}

	return &processor, nil
}

func (i *ProcessorRepo) Create(ctx context.Context, processor *models.Processor) error {
	mutateFn := func(processor *models.Processor) error {
		return sqlDB.WithContext(ctx).Create(processor).Error
	}

	cl := cache.NewClient[*models.Processor](0)
	invKeys := []string{WorkspaceProcessorsKey(processor.WorkspaceID)}
	key := ProcessorKey(processor.ID)

	if err := cl.Mutate(ctx, key, invKeys, processor, mutateFn, 0); err != nil {
		return err
	}

	return nil
}

func (i *ProcessorRepo) Patch(ctx context.Context, workspaceID, processorID string, patch map[string]interface{}) error {
	if err := sqlDB.WithContext(ctx).Model(&models.Processor{}).Where("id = ?", processorID).Updates(patch).Error; err != nil {
		return dbutils.DBError(err)
	}
	invKeys := []string{WorkspaceProcessorsKey(workspaceID), ProcessorKey(processorID)}
	cache.Invalidate(ctx, invKeys...)
	return nil
}

func (i *ProcessorRepo) Delete(ctx context.Context, workspaceID, processorID string) error {
	if err := sqlDB.WithContext(ctx).Delete(&models.Processor{}, "id = ?", processorID).Error; err != nil {
		return dbutils.DBError(err)
	}
	invKeys := []string{WorkspaceProcessorsKey(workspaceID), ProcessorKey(processorID)}
	cache.Invalidate(ctx, invKeys...)
	return nil
}

func (i *ProcessorRepo) SaveWorkflow(ctx context.Context, processorID string, workflow *models.Workflow) error {
	patch := map[string]interface{}{
		"statement": workflow.Root,
		"variables": workflow.Variables,
	}
	if err := sqlDB.WithContext(ctx).Model(&models.Processor{}).Where("id = ?", processorID).Updates(patch).Error; err != nil {
		return dbutils.DBError(err)
	}

	invKeys := []string{ProcessorKey(processorID)}
	cache.Invalidate(ctx, invKeys...)
	return nil
}

func WorkspaceProcessorsKey(workspaceID string) string {
	return "workspace:" + workspaceID + ":processors"
}

func ProcessorKey(processorID string) string {
	return "processor:" + processorID
}
