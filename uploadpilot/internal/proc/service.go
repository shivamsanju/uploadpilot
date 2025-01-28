package proc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
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
	return s.wsRepo.GetAll(ctx, workspaceID)
}

func (s *ProcessorService) GetProcessor(ctx context.Context, workspaceID string, processorID string) (*models.Processor, error) {
	return s.wsRepo.Get(ctx, workspaceID, processorID)
}

func (s *ProcessorService) CreateProcessor(ctx context.Context, workspaceID string, processor *models.Processor) error {
	return s.wsRepo.Create(ctx, workspaceID, processor)
}

func (s *ProcessorService) UpdateTasks(ctx context.Context, workspaceID string, processorID string, tasks *models.ProcTaskCanvas) error {
	patch := bson.M{"tasks": tasks}
	return s.wsRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *ProcessorService) DeleteProcessor(ctx context.Context, workspaceID string, processorID string) error {
	return s.wsRepo.Delete(ctx, workspaceID, processorID)
}

func (s *ProcessorService) EnableDisableProcessor(ctx context.Context, workspaceID string, processorID string, enabled bool) error {
	patch := bson.M{"enabled": enabled}
	return s.wsRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *ProcessorService) EditNameAndTrigger(ctx context.Context, workspaceID string, processorID string, update *dto.EditProcRequest) error {
	patch := bson.M{"name": update.Name, "triggers": update.Triggers}
	return s.wsRepo.Patch(ctx, workspaceID, processorID, patch)
}
