package processor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/uploadpilot/go-core/dsl"
	"github.com/uploadpilot/uploadpilot/uploader/internal/dto"
	"github.com/uploadpilot/uploadpilot/uploader/internal/infra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"gopkg.in/yaml.v3"
)

type Service struct {
	processorRepo *repo.ProcessorRepo
}

func NewProcessorService(processorRepo *repo.ProcessorRepo) *Service {
	return &Service{
		processorRepo: processorRepo,
	}
}

func (s *Service) GetProcessors(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.processorRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *Service) TriggerWorkflow(ctx context.Context, workspaceID string, processor *models.Processor) (*dto.TriggerWorkflowResp, error) {
	var dslWorkflow dsl.Workflow
	if err := yaml.Unmarshal([]byte(processor.Workflow), &dslWorkflow); err != nil {
		log.Fatalln("failed to unmarshal dsl config", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "dsl",
		TypedSearchAttributes: temporal.NewSearchAttributes(
			temporal.NewSearchAttributeKeyKeyword("processorId").ValueSet(processor.ID),
		),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    processor.MaxRetries,
			InitialInterval:    time.Duration(processor.RetryInitialIntervalS) * time.Second,
			BackoffCoefficient: processor.RetryBackoffCoefficient,
			MaximumInterval:    time.Duration(processor.RetryMaxIntervalS) * time.Second,
		},
		WorkflowExecutionTimeout: time.Duration(processor.WorkflowRunTimeoutS) * time.Second,
		WorkflowRunTimeout:       time.Duration(processor.WorkflowExecutionTimeoutS) * time.Second,
	}

	we, err := infra.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		infra.Log.Errorf("Unable to execute workflow", err)
		return nil, err
	}
	infra.Log.Infof("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	return &dto.TriggerWorkflowResp{WorkflowID: we.GetID(), RunID: we.GetRunID()}, nil
}
