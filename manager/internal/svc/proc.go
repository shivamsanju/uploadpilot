package svc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/validations"
	"github.com/uploadpilot/uploadpilot/common/pkg/validations/jsonschema"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
	"gopkg.in/yaml.v3"
)

type ProcessorService struct {
	procRepo *repo.ProcessorRepo
}

func NewProcessorService(procRepo *repo.ProcessorRepo) *ProcessorService {
	return &ProcessorService{
		procRepo: procRepo,
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

	return s.procRepo.Create(ctx, processor)
}

func (s *ProcessorService) UpdateWorkflow(ctx context.Context, workspaceID, processorID string, workflow string) error {
	// TODO: Validate workflow
	//TODO: Validate tasks
	var json map[string]interface{}
	if err := yaml.Unmarshal([]byte(workflow), &json); err != nil {
		infra.Log.Errorf("failed to unmarshal workflow: %s", err.Error())
		return err
	}

	if err := validations.ValidateJSONSchema(jsonschema.WorkflowSchema, json); err != nil {
		return err
	}

	return s.procRepo.SaveWorkflow(ctx, workspaceID, processorID, workflow)
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
