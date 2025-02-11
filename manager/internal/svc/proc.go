package svc

import (
	"context"
	"fmt"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type ProcessorService struct {
	procRepo *db.ProcessorRepo
	taskRepo *db.TaskRepo
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{
		procRepo: db.NewProcessorRepo(),
		taskRepo: db.NewTaskRepo(),
	}
}

func (s *ProcessorService) GetAllProcessorsInWorkspace(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.procRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *ProcessorService) GetProcessor(ctx context.Context, processorID string) (*models.Processor, error) {
	processor, err := s.procRepo.Get(ctx, processorID)
	if err != nil {
		return nil, err
	}
	return processor, nil
}

func (s *ProcessorService) CreateProcessor(ctx context.Context, workspaceID string, processor *models.Processor) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}

	processor.CreatedBy = user.UserID
	processor.UpdatedBy = user.UserID
	processor.WorkspaceID = workspaceID
	processor.Tasks = []models.Task{}

	return s.procRepo.Create(ctx, processor)
}

func (s *ProcessorService) GetTasks(ctx context.Context, workspaceID, processorID string) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetAll(ctx, processorID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *ProcessorService) SaveTasks(ctx context.Context, workspaceID, processorID string, tasks []models.Task) error {
	//TODO: Validate tasks
	return s.taskRepo.SaveTasks(ctx, processorID, tasks)
}

func (s *ProcessorService) UpdateWorkflow(ctx context.Context, workspaceID, processorID string, workflow *models.Workflow) error {
	infra.Log.Infof("WFLOW %+v", workflow)
	if workflow.Root.Activity == nil && workflow.Root.Sequence == nil && workflow.Root.Parallel == nil {
		return fmt.Errorf("atleast one statement is required")
	}
	return s.procRepo.SaveWorkflow(ctx, processorID, workflow)
}

func (s *ProcessorService) DeleteProcessor(ctx context.Context, workspaceID, processorID string) error {
	return s.procRepo.Delete(ctx, workspaceID, processorID)
}

func (s *ProcessorService) EnableDisableProcessor(ctx context.Context, workspaceID, processorID string, enabled bool) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"enabled": enabled}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *ProcessorService) EditNameAndTrigger(ctx context.Context, workspaceID, processorID string, update *dto.EditProcRequest) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"name": update.Name, "triggers": update.Triggers}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}
