package proc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	return s.wsRepo.GetAll(ctx, workspaceID)
}

func (s *ProcessorService) GetProcessor(ctx context.Context, processorID string) (*models.Processor, error) {
	return s.wsRepo.Get(ctx, processorID)
}

func (s *ProcessorService) CreateProcessor(ctx context.Context, workspaceID string, processor *models.Processor) error {
	wsObjID, err := primitive.ObjectIDFromHex(workspaceID)
	if err != nil {
		return err
	}
	processor.WorkspaceID = wsObjID
	processor.Tasks = &models.ProcTaskCanvas{
		Nodes: []models.ProcTask{
			{
				ID:        primitive.NewObjectID().Hex(),
				Key:       TriggerTaskKey,
				Type:      "baseNode",
				Deletable: false,
				Data: models.JSON{
					"label":       "Trigger",
					"description": "Trigger the processor to start processing the files",
				},
				Position: models.JSON{
					"x": 0,
					"y": 0,
				},
				Measured: models.JSON{},
			},
		},
		Edges: []models.ProcTaskEdge{},
	}
	return s.wsRepo.Create(ctx, processor)
}

func (s *ProcessorService) UpdateTasks(ctx context.Context, processorID string, tasks *models.ProcTaskCanvas) error {
	patch := bson.M{"tasks": tasks}
	return s.wsRepo.Patch(ctx, processorID, patch)
}

func (s *ProcessorService) DeleteProcessor(ctx context.Context, processorID string) error {
	return s.wsRepo.Delete(ctx, processorID)
}

func (s *ProcessorService) EnableDisableProcessor(ctx context.Context, processorID string, enabled bool) error {
	patch := bson.M{"enabled": enabled}
	return s.wsRepo.Patch(ctx, processorID, patch)
}

func (s *ProcessorService) EditNameAndTrigger(ctx context.Context, processorID string, update *dto.EditProcRequest) error {
	patch := bson.M{"name": update.Name, "triggers": update.Triggers}
	return s.wsRepo.Patch(ctx, processorID, patch)
}
