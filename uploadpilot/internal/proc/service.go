package proc

import (
	"context"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/db/types"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProcessorService struct {
	wsRepo *db.ProcessorRepo
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{
		wsRepo: db.NewProcessorRepo(),
	}
}

func (s *ProcessorService) GetAllProcessorsInWorkspace(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.wsRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *ProcessorService) GetProcessor(ctx context.Context, processorID string) (*models.Processor, error) {
	processor, err := s.wsRepo.Get(ctx, processorID)
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
	processor.Canvas = types.JSONB{
		"nodes": []types.JSONB{{
			"id":        uuid.NewString(),
			"key":       TriggerTaskKey,
			"type":      "baseNode",
			"deletable": false,
			"data": types.EncryptedJSONB{
				"label":       "Trigger",
				"description": "Trigger the processor to start processing the files",
			},
			"position": types.JSONB{
				"x": 0,
				"y": 0,
			},
			"measured": types.JSONB{},
		}},
		"edges": []types.JSONB{},
	}
	return s.wsRepo.Create(ctx, processor)
}

func (s *ProcessorService) UpdateTasks(ctx context.Context, processorID string, tasks *types.JSONB) error {
	patch := bson.M{"tasks": tasks}
	return s.wsRepo.Patch(ctx, processorID, patch)
}

func (s *ProcessorService) DeleteProcessor(ctx context.Context, workspaceID, processorID string) error {
	return s.wsRepo.Delete(ctx, processorID)
}

func (s *ProcessorService) EnableDisableProcessor(ctx context.Context, workspaceID, processorID string, enabled bool) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"enabled": enabled}
	patch["updated_by"] = user.UserID
	return s.wsRepo.Patch(ctx, processorID, patch)
}

func (s *ProcessorService) EditNameAndTrigger(ctx context.Context, workspaceID, processorID string, update *dto.EditProcRequest) error {
	user, err := utils.GetUserDetailsFromContext(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"name": update.Name, "triggers": update.Triggers}
	patch["updated_by"] = user.UserID
	return s.wsRepo.Patch(ctx, processorID, patch)
}
