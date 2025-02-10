package db

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/cache"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type TaskRepo struct {
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{}
}

func (i *TaskRepo) GetAll(ctx context.Context, processorID string) ([]models.Task, error) {
	dbFetchFn := func(tasks *[]models.Task) error {
		return sqlDB.WithContext(ctx).Where("processor_id = ?", processorID).Find(tasks).Error
	}

	var tasks []models.Task
	cl := cache.NewClient[*[]models.Task](0)
	key := ProcessorTaskKey(processorID)

	if err := cl.Query(ctx, key, &tasks, dbFetchFn); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (i *TaskRepo) SaveTasks(ctx context.Context, processorID string, tasks []models.Task) error {
	mutateFn := func(tasks *[]models.Task) error {
		return sqlDB.WithContext(ctx).Save(tasks).Error
	}

	cl := cache.NewClient[*[]models.Task](0)
	invKeys := []string{ProcessorTaskKey(processorID)}
	key := ProcessorTaskKey(processorID)

	if err := cl.Mutate(ctx, key, invKeys, &tasks, mutateFn, 0); err != nil {
		return err
	}

	return nil
}

func ProcessorTaskKey(processorID string) string {
	return "processor:" + processorID + ":tasks"
}

func TaskKey(taskID string) string {
	return "task:" + taskID
}
