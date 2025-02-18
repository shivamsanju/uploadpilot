package svc

import (
	"context"
	"fmt"

	"github.com/uploadpilot/uploadpilot/go-core/common/tasks"
	"github.com/uploadpilot/uploadpilot/go-core/common/validator"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/uploadpilot/go-core/dsl"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
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

	if err := validator.ValidateJSONSchema(dsl.DSLSchema, json); err != nil {
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

func (s *ProcessorService) GetAllTasks(ctx context.Context) []tasks.Task {
	var tsks []tasks.Task
	for _, task := range tasks.TaskCatalog {
		tsks = append(tsks, *task)
	}
	return tsks
}

func (s *ProcessorService) GetWorkflowRuns(ctx context.Context, processorID string) ([]dto.WorkflowRun, error) {
	result, err := infra.TemporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		Query: "processorId = '" + processorID + "'",
	})
	if err != nil {
		return nil, err
	}
	var runs []dto.WorkflowRun
	for _, run := range result.Executions {
		runs = append(runs, dto.WorkflowRun{
			WorkflowID:      run.Execution.WorkflowId,
			RunID:           run.Execution.RunId,
			StartTime:       run.StartTime.AsTime(),
			EndTime:         run.CloseTime.AsTime(),
			DurationSeconds: run.CloseTime.Seconds - run.StartTime.Seconds,
			Status:          run.Status.Enum().String(),
		})
	}
	return runs, nil
}

func (s *ProcessorService) GetWorkflowHistory(ctx context.Context, workflowID string, runID string) ([]dto.WorkflowRunLogs, error) {
	infra.Log.Infof("Getting workflow history for workflowID: %s, runID: %s", workflowID, runID)
	iter := infra.TemporalClient.GetWorkflowHistory(context.Background(), workflowID, runID, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	var logs []dto.WorkflowRunLogs
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			return nil, err
		}
		eventType := event.GetEventType()

		log := dto.WorkflowRunLogs{
			Timestamp: event.EventTime.AsTime(),
			EventType: event.GetEventType().String(),
		}

		switch eventType {
		case enums.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED:
			attributes := event.GetWorkflowExecutionStartedEventAttributes()
			payloads := attributes.Input.Payloads
			payloadStr := ""
			for _, payload := range payloads {
				payloadStr += string(payload.Data)
			}
			log.Details = fmt.Sprintf("Input: %s", payloadStr)
		case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			attributes := event.GetActivityTaskScheduledEventAttributes()
			payloads := attributes.Input.Payloads
			payloadStr := ""
			for _, payload := range payloads {
				payloadStr += string(payload.Data)
			}
			log.Details = fmt.Sprintf("ActivityType: %s. Input: %s", attributes.ActivityType.Name, payloadStr)

		case enums.EVENT_TYPE_ACTIVITY_TASK_STARTED:
			attributes := event.GetActivityTaskStartedEventAttributes()
			log.Details = fmt.Sprintf("Attempt: %d", attributes.Attempt)

		case enums.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
			attributes := event.GetActivityTaskCompletedEventAttributes()
			log.Details = fmt.Sprintf("Result: %s", attributes.Result)

		case enums.EVENT_TYPE_ACTIVITY_TASK_FAILED:
			attributes := event.GetActivityTaskFailedEventAttributes()
			log.Details = fmt.Sprintf("Cause: %s, Message: %s", attributes.Failure.Cause, attributes.Failure.Message)

		case enums.EVENT_TYPE_ACTIVITY_TASK_TIMED_OUT:
			attributes := event.GetActivityTaskTimedOutEventAttributes()
			log.Details = fmt.Sprintf("Cause: %s, Message: %s", attributes.Failure.Cause, attributes.Failure.Message)

		case enums.EVENT_TYPE_WORKFLOW_EXECUTION_FAILED:
			attributes := event.GetWorkflowExecutionFailedEventAttributes()
			log.Details = fmt.Sprintf("Cause: %s, Message: %s", attributes.Failure.Cause, attributes.Failure.Message)

		case enums.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED:
			attributes := event.GetWorkflowExecutionCompletedEventAttributes()
			payloads := attributes.Result.Payloads
			payloadStr := ""
			for _, payload := range payloads {
				payloadStr += string(payload.Data)
			}
			log.Details = fmt.Sprintf("Result: %s", payloadStr)

		default:
			log.Details = ""
		}

		logs = append(logs, log)
	}

	return logs, nil
}
