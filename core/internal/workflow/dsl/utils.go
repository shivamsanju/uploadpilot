package dsl

import (
	"encoding/json"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func addWorkflowIdentifiersToBindings(ctx workflow.Context, bindings map[string]any, dslWorkflow Workflow) {
	bindings["workspace_id"] = dslWorkflow.WorkspaceID
	bindings["upload_id"] = dslWorkflow.UploadID
	bindings["processor_id"] = dslWorkflow.ProcessorID
	bindings["file_name"] = dslWorkflow.FileName
	bindings["content_type"] = dslWorkflow.ContentType
	bindings["workflow_id"] = workflow.GetInfo(ctx).WorkflowExecution.ID
	bindings["run_id"] = workflow.GetInfo(ctx).WorkflowExecution.RunID
}

func makeInput(argMap map[string]any, bindings map[string]any, activityKey string, saveOutput *bool, inputActivityKey *string) (string, error) {
	for argument, value := range argMap {
		val, ok := value.(string)
		if ok && val[0] == '$' {
			varName := val[1:]
			bindings[fmt.Sprintf("%s.%s", activityKey, argument)] = bindings[varName]
		} else {
			bindings[fmt.Sprintf("%s.%s", activityKey, argument)] = value
		}
	}

	bindings["current_activity_key"] = activityKey
	if saveOutput != nil {
		bindings[fmt.Sprintf("%s.save_output", activityKey)] = *saveOutput
	} else {
		bindings[fmt.Sprintf("%s.save_output", activityKey)] = false
	}
	if inputActivityKey != nil {
		bindings[fmt.Sprintf("%s.input", activityKey)] = *inputActivityKey
	} else {
		bindings[fmt.Sprintf("%s.input", activityKey)] = ""
	}

	argsbytes, err := json.Marshal(bindings)
	if err != nil {
		return "", err
	}

	return string(argsbytes), nil
}

func saveOutput(result map[string]any, bindings map[string]any, activityKey string) error {
	for key, value := range result {
		bindings[fmt.Sprintf("%s.%s", activityKey, key)] = value
	}
	return nil
}

func runPostProcessingActivity(ctx workflow.Context, bindings map[string]any) error {
	ao := workflow.ActivityOptions{}
	ao.RetryPolicy = &temporal.RetryPolicy{
		MaximumAttempts:    1,
		InitialInterval:    0,
		BackoffCoefficient: 2,
		MaximumInterval:    1 * time.Minute,
	}
	ao.ScheduleToStartTimeout = 24 * time.Hour
	ao.StartToCloseTimeout = 24 * time.Hour
	ao.ScheduleToCloseTimeout = 24 * time.Hour

	ctx = workflow.WithActivityOptions(ctx, ao)

	bindingsB, err := json.Marshal(bindings)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, "Executor", "PostProcessingV1", string(bindingsB)).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
