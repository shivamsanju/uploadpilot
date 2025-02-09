package proc

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
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
