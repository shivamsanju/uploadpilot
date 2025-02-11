package proc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/dsl"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/uploader/internal/dto"
	"go.temporal.io/sdk/client"
	"gopkg.in/yaml.v3"
)

type ProcessorService struct {
	pRepo *db.ProcessorRepo
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{
		pRepo: db.NewProcessorRepo(),
	}
}

func (s *ProcessorService) GetProcessors(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.pRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *ProcessorService) TriggerWorkflow(ctx context.Context, yamlContent string) (*dto.TriggerWorkflowResp, error) {
	var dslWorkflow dsl.Workflow
	if err := yaml.Unmarshal([]byte(yamlContent), &dslWorkflow); err != nil {
		log.Fatalln("failed to unmarshal dsl config", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "dsl",
	}

	we, err := infra.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		infra.Log.Errorf("Unable to execute workflow", err)
		return nil, err
	}
	infra.Log.Infof("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	return &dto.TriggerWorkflowResp{WorkflowID: we.GetID(), RunID: we.GetRunID()}, nil
}
