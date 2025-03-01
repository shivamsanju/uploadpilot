package processor

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/phuslu/log"
	"github.com/uploadpilot/go-common/workflow/catalog"
	"github.com/uploadpilot/go-common/workflow/dsl"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/db/repo"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/manager/internal/svc/processor/templates"
	"github.com/uploadpilot/manager/internal/utils"
	"github.com/uploadpilot/manager/internal/validator"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"gopkg.in/yaml.v3"
)

type Service struct {
	procRepo *repo.ProcessorRepo
}

func NewService(procRepo *repo.ProcessorRepo) *Service {
	return &Service{
		procRepo: procRepo,
	}
}

func (s *Service) GetAllProcessorsInWorkspace(ctx context.Context, workspaceID string) ([]models.Processor, error) {
	processors, err := s.procRepo.GetAll(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return processors, nil
}

func (s *Service) GetProcessor(ctx context.Context, processorID string) (*models.Processor, error) {
	processor, err := s.procRepo.Get(ctx, processorID)
	if err != nil {
		return nil, err
	}
	return processor, nil
}

func (s *Service) GetTemplates(ctx context.Context) []dto.ProcessorTemplate {
	return templates.ProcessorTemplates
}

func (s *Service) CreateProcessor(ctx context.Context, workspaceID string, processor *models.Processor, templateKey *string) error {
	user, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	if templateKey == nil || *templateKey == "" {
		templateKey = new(string)
		*templateKey = "sample"
	}
	sampleWflow, err := os.ReadFile("./internal/svc/processor/templates/" + *templateKey + ".yaml")
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

func (s *Service) UpdateWorkflow(ctx context.Context, workspaceID, processorID string, workflow string) error {
	// TODO: Validate workflow
	//TODO: Validate tasks
	var json map[string]interface{}
	if err := yaml.Unmarshal([]byte(workflow), &json); err != nil {
		log.Error().Msgf("failed to unmarshal workflow: %s", err.Error())
		return err
	}

	if err := validator.ValidateJSONSchema(dsl.DSLSchema, json); err != nil {
		return err
	}

	return s.procRepo.SaveWorkflow(ctx, workspaceID, processorID, workflow)
}

func (s *Service) DeleteProcessor(ctx context.Context, workspaceID, processorID string) error {
	return s.procRepo.Delete(ctx, workspaceID, processorID)
}

func (s *Service) EnableDisableProcessor(ctx context.Context, workspaceID, processorID string, enabled bool) error {
	user, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"enabled": enabled}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *Service) EditNameAndTrigger(ctx context.Context, workspaceID, processorID string, update *dto.EditProcRequest) error {
	user, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}
	patch := map[string]interface{}{"name": update.Name, "triggers": update.Triggers}
	patch["updated_by"] = user.UserID
	return s.procRepo.Patch(ctx, workspaceID, processorID, patch)
}

func (s *Service) GetAllTasks(ctx context.Context) []catalog.ActivityMetadata {
	var tsks []catalog.ActivityMetadata
	for _, task := range catalog.ActivityCatalog {
		tsks = append(tsks, *task)
	}
	return tsks
}

func (s *Service) TriggerWorkflows(ctx context.Context, workspaceID string, upload *models.Upload) error {
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
				doTrigger = true
			}
			if !doTrigger {
				continue
			}
			_, err := s.TriggerWorkflow(ctx, upload, processor.Workflow, workspaceID, processor.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) TriggerWorkflow(ctx context.Context, upload *models.Upload, yamlContent, workspaceID, processorID string) (*dto.TriggerWorkflowResp, error) {
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

	we, err := infra.TemporalClient.ExecuteWorkflow(context.Background(), workflowOptions, dsl.SimpleDSLWorkflow, dslWorkflow)
	if err != nil {
		log.Error().Err(err).Msg("failed to start workflow")
		return nil, err
	}
	return &dto.TriggerWorkflowResp{WorkflowID: we.GetID(), RunID: we.GetRunID()}, nil
}

func (s *Service) GetWorkflowRuns(ctx context.Context, processorID string) ([]dto.WorkflowRun, error) {
	result, err := infra.TemporalClient.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
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

func (s *Service) GetWorkflowHistory(ctx context.Context, workflowID string, runID string) ([]dto.WorkflowRunLogs, error) {
	log.Info().Msgf("Getting workflow history for workflowID: %s, runID: %s", workflowID, runID)
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
