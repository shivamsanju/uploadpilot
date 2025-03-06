package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/uploadpilot/core/internal/activities/catalog"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/templates"
	"github.com/uploadpilot/core/pkg/dsl"
	"github.com/uploadpilot/core/pkg/validator"
	"github.com/uploadpilot/core/web/webutils"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"gopkg.in/yaml.v3"
)

type ProcessorService struct {
	procRepo       *repo.ProcessorRepo
	validator      *validator.Validator
	temporalClient client.Client
}

func NewProcessorService(procRepo *repo.ProcessorRepo, temporalClient client.Client) *ProcessorService {
	return &ProcessorService{
		procRepo:       procRepo,
		validator:      validator.NewValidator(),
		temporalClient: temporalClient,
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

func (s *ProcessorService) GetTemplates(ctx context.Context) []dto.ProcessorTemplate {
	return templates.ProcessorTemplates
}

func (s *ProcessorService) CreateProcessor(ctx context.Context, workspaceID string, processor *models.Processor, templateKey *string) error {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if templateKey == nil || *templateKey == "" {
		templateKey = new(string)
		*templateKey = "sample"
	}
	sampleWflow, err := os.ReadFile("./internal/templates/" + *templateKey + ".yaml")
	if err != nil {
		return err
	}
	wfData := string(sampleWflow)

	processor.CreatedBy = user.UserID
	processor.UpdatedBy = user.UserID
	processor.WorkspaceID = workspaceID
	processor.Workflow = wfData

	return s.procRepo.Create(ctx, processor)
}

func (s *ProcessorService) UpdateWorkflow(ctx context.Context, workspaceID, processorID string, workflow string) error {
	// TODO: Validate workflow
	//TODO: Validate tasks
	var json map[string]interface{}
	if err := yaml.Unmarshal([]byte(workflow), &json); err != nil {
		log.Error().Msgf("failed to unmarshal workflow: %s", err.Error())
		return err
	}

	if err := s.validator.ValidateJSONSchema(dsl.DSLSchema, json); err != nil {
		return err
	}

	return s.procRepo.SaveWorkflow(ctx, workspaceID, processorID, workflow)
}

func (s *ProcessorService) DeleteProcessor(ctx context.Context, workspaceID, processorID string) error {
	return s.procRepo.Delete(ctx, workspaceID, processorID)
}

func (s *ProcessorService) EnableDisableProcessor(ctx context.Context, workspaceID, processorID string, enabled bool) error {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"enabled": enabled}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *ProcessorService) EditNameAndTrigger(ctx context.Context, workspaceID, processorID string, update *dto.EditProcRequest) error {
	user, err := webutils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"name": update.Name, "triggers": update.Triggers}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *ProcessorService) GetAllActivities(ctx context.Context) []catalog.ActivityMetadata {
	var tsks []catalog.ActivityMetadata
	for _, task := range catalog.ActivityCatalog {
		tsks = append(tsks, *task)
	}
	return tsks
}

func (s *ProcessorService) TriggerWorkflows(ctx context.Context, workspaceID string, upload *models.Upload) error {
	processors, err := s.GetAllProcessorsInWorkspace(ctx, workspaceID)
	if err != nil {
		return err
	}
	for _, processor := range processors {
		if processor.Enabled {
			var doTrigger bool = false
			if len(processor.Triggers) != 0 {
				for _, trigger := range processor.Triggers {
					if trigger == upload.FileType {
						doTrigger = true
						break
					}
				}
			} else {
				doTrigger = false
			}
			if !doTrigger {
				continue
			}
			log.Debug().Msgf("Triggering processor: %s, Workflow: %s", processor.Name, processor.Workflow)
			_, err := s.TriggerWorkflow(ctx, upload, processor.Workflow, workspaceID, processor.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ProcessorService) TriggerWorkflow(ctx context.Context, upload *models.Upload, yamlContent, workspaceID, processorID string) (*dto.TriggerWorkflowResp, error) {
	var dslWorkflow dsl.Workflow
	if err := yaml.Unmarshal([]byte(yamlContent), &dslWorkflow); err != nil {
		return nil, err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: "queue1",
		TypedSearchAttributes: temporal.NewSearchAttributes(
			temporal.NewSearchAttributeKeyKeyword("processorId").ValueSet(processorID),
		),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
		Memo: map[string]interface{}{
			"uploadId":    upload.ID,
			"workspaceId": workspaceID,
			"fileType":    upload.FileType,
			"fileName":    upload.FileName,
		},
	}

	dslWorkflow.WorkspaceID = workspaceID
	dslWorkflow.UploadID = upload.ID
	dslWorkflow.ProcessorID = processorID
	dslWorkflow.UploadFileName = upload.FileName
	dslWorkflow.UploadFileType = upload.FileType

	fmt.Println("workflowOptions", workflowOptions)
	we, err := s.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		log.Error().Err(err).Msg("failed to start workflow")
		return nil, err
	}
	return &dto.TriggerWorkflowResp{WorkflowID: we.GetID(), RunID: we.GetRunID()}, nil
}

func (s *ProcessorService) GetWorkflowRuns(ctx context.Context, processorID string) ([]dto.WorkflowRun, error) {
	result, err := s.temporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		Query: "processorId = '" + processorID + "'",
	})
	if err != nil {
		return nil, err
	}
	var runs []dto.WorkflowRun
	for _, run := range result.Executions {

		r := dto.WorkflowRun{
			ID:         run.Execution.RunId,
			WorkflowID: run.Execution.WorkflowId,
			RunID:      run.Execution.RunId,
			StartTime:  run.StartTime.AsTime(),
			Status:     run.Status.Enum().String(),
		}
		if workspaceIDField, ok := run.Memo.Fields["workspaceId"]; ok && workspaceIDField != nil {
			r.WorkspaceID = strings.Trim(string(workspaceIDField.GetData()), "\"")
		} else {
			return nil, fmt.Errorf("workspaceId not found in run")
		}

		if uploadIDField, ok := run.Memo.Fields["uploadId"]; ok && uploadIDField != nil {
			r.UploadID = strings.Trim(string(uploadIDField.GetData()), "\"")
		} else {
			return nil, fmt.Errorf("uploadId not found in run")
		}

		if run.CloseTime != nil {
			r.EndTime = run.CloseTime.AsTime()
			r.WorkflowTimeMillis = run.CloseTime.AsTime().UnixMilli() - run.StartTime.AsTime().UnixMilli()
			r.ExecutionTimeMilis = run.ExecutionDuration.AsDuration().Milliseconds()
		} else {
			r.WorkflowTimeMillis = time.Now().Unix() - run.StartTime.Seconds
		}
		runs = append(runs, r)

	}
	return runs, nil
}

func (s *ProcessorService) GetWorkflowHistory(ctx context.Context, workflowID string, runID string) ([]dto.WorkflowRunLogs, error) {
	log.Info().Msgf("Getting workflow history for workflowID: %s, runID: %s", workflowID, runID)
	var logs []dto.WorkflowRunLogs

	iter := s.temporalClient.GetWorkflowHistory(context.Background(), workflowID, runID, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
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
			log.Details = fmt.Sprintf("Timeout: %s", attributes.WorkflowRunTimeout)
		case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			attributes := event.GetActivityTaskScheduledEventAttributes()
			payloads := attributes.Input.Payloads
			payloadStr := ""
			for _, payload := range payloads {
				payloadStr += string(payload.Data)
			}
			log.Details = fmt.Sprintf("ActivityType: %s,  Input: %s", attributes.ActivityType.Name, payloadStr)

		case enums.EVENT_TYPE_ACTIVITY_TASK_STARTED:
			attributes := event.GetActivityTaskStartedEventAttributes()
			log.Details = fmt.Sprintf("Attempt: %d", attributes.Attempt)

		case enums.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
			attributes := event.GetActivityTaskCompletedEventAttributes()
			log.Details = fmt.Sprintf("Result: %s", attributes.Result)

		case enums.EVENT_TYPE_ACTIVITY_TASK_FAILED:
			attributes := event.GetActivityTaskFailedEventAttributes()
			log.Details = fmt.Sprintf("Message: %s", attributes.Failure.Message)

		case enums.EVENT_TYPE_ACTIVITY_TASK_TIMED_OUT:
			attributes := event.GetActivityTaskTimedOutEventAttributes()
			log.Details = fmt.Sprintf("Message: %s", attributes.Failure.Message)

		case enums.EVENT_TYPE_WORKFLOW_EXECUTION_FAILED:
			attributes := event.GetWorkflowExecutionFailedEventAttributes()
			log.Details = fmt.Sprintf("Message: %s", attributes.Failure.Message)

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
