package dsl

import (
	"encoding/json"
	"fmt"
	"time"

	"maps"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	WorkflowCtxKey string

	Workflow struct {
		WorkspaceID       string         `json:"workspaceId"`
		UploadID          string         `json:"uploadId"`
		ProcessorID       string         `json:"processorId"`
		FileName          string         `json:"fileName"`
		ContentType       string         `json:"contentType"`
		Variables         map[string]any `json:"variables" yaml:"variables"`
		Root              Statement      `json:"root" yaml:"root"`
		OnWorkflowSuccess *Statement     `json:"on_workflow_success" yaml:"on_workflow_success"`
		OnWorkflowFailure *Statement     `json:"on_workflow_failure" yaml:"on_workflow_failure"`
	}

	Statement struct {
		Activity *ActivityInvocation `json:"activity,omitempty" yaml:"activity,omitempty"`
		Sequence *Sequence           `json:"sequence,omitempty" yaml:"sequence,omitempty"`
		Parallel *Parallel           `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	Sequence struct {
		Elements []*Statement `json:"elements" yaml:"elements"`
	}

	Parallel struct {
		Branches []*Statement `json:"branches" yaml:"branches"`
	}

	ActivityInvocation struct {
		Key                           string         `json:"key" yaml:"key"`
		Uses                          string         `json:"uses" yaml:"uses"`
		With                          map[string]any `json:"with,omitempty" yaml:"with,omitempty"`
		Input                         *string        `json:"input,omitempty" yaml:"input,omitempty"`
		SaveOutput                    *bool          `json:"save_output,omitempty" yaml:"save_output,omitempty"`
		ScheduleToCloseTimeoutSeconds *int64         `json:"schedule_to_close_timeout_seconds,omitempty" yaml:"schedule_to_close_timeout_seconds,omitempty"`
		ScheduleToStartTimeoutSeconds *int64         `json:"schedule_to_start_timeout_seconds,omitempty" yaml:"schedule_to_start_timeout_seconds,omitempty"`
		StartToCloseTimeoutSeconds    *int64         `json:"start_to_close_timeout_seconds,omitempty" yaml:"start_to_close_timeout_seconds,omitempty"`
		MaxRetries                    *int32         `json:"max_retries,omitempty" yaml:"max_retries,omitempty"`
		RetryBackoffCoefficient       *float64       `json:"retry_backoff_coefficient,omitempty" yaml:"retry_backoff_coefficient,omitempty"`
		RetryMaxIntervalSeconds       *int64         `json:"retry_max_interval_seconds,omitempty" yaml:"retry_max_interval_seconds,omitempty"`
		RetryInitialIntervalSeconds   *int64         `json:"retry_initial_interval_seconds,omitempty" yaml:"retry_initial_interval_seconds,omitempty"`
		OnSuccess                     *Statement     `json:"on_success,omitempty" yaml:"on_success,omitempty"`
		OnError                       *Statement     `json:"on_error,omitempty" yaml:"on_error,omitempty"`
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]any) error
	}
)

func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	logger := workflow.GetLogger(ctx)
	bindings := make(map[string]any)
	maps.Copy(bindings, dslWorkflow.Variables)

	// adds workspace_id, upload_id, run_id etc
	addWorkflowIdentifiersToBindings(ctx, bindings, dslWorkflow)

	workflowErr := dslWorkflow.Root.execute(ctx, bindings)

	// runs the post processing activity in any case
	// TODO: handle what if it fails
	if e := runPostProcessingActivity(ctx, bindings); e != nil {
		logger.Error("Error in post processing: ", e)
		if workflowErr != nil {
			workflowErr = fmt.Errorf("failed to run post processing. original error: %w", workflowErr)
		} else {
			workflowErr = fmt.Errorf("failed to run post processing")
		}
	}

	if workflowErr != nil {
		logger.Error("DSL Workflow failed: ", workflowErr)
		bindings["workflow_error"] = workflowErr

		if dslWorkflow.OnWorkflowFailure != nil {
			onFailureErr := dslWorkflow.OnWorkflowFailure.execute(ctx, bindings)
			if onFailureErr != nil {
				workflowErr = fmt.Errorf("failed to run on_workflow_failure: %w. original error: %w", onFailureErr, workflowErr)
			}
		}
		return nil, workflowErr
	}

	if dslWorkflow.OnWorkflowSuccess != nil {
		onSuccessErr := dslWorkflow.OnWorkflowSuccess.execute(ctx, bindings)
		if onSuccessErr != nil {
			return nil, fmt.Errorf("failed to run on_workflow_success: %w", onSuccessErr)
		}
	}

	logger.Info("DSL Workflow completed.")
	return nil, nil
}

func (b *Statement) execute(ctx workflow.Context, bindings map[string]any) error {
	if b.Parallel != nil {
		return b.Parallel.execute(ctx, bindings)
	}
	if b.Sequence != nil {
		return b.Sequence.execute(ctx, bindings)
	}
	if b.Activity != nil {
		return b.Activity.execute(ctx, bindings)
	}
	return nil
}

func (a *ActivityInvocation) execute(ctx workflow.Context, bindings map[string]any) error {
	ao := workflow.ActivityOptions{}
	if a.StartToCloseTimeoutSeconds != nil && *a.StartToCloseTimeoutSeconds != 0 {
		ao.StartToCloseTimeout = time.Duration(*a.StartToCloseTimeoutSeconds) * time.Second
	} else {
		ao.StartToCloseTimeout = 24 * time.Hour
	}

	if a.ScheduleToCloseTimeoutSeconds != nil {
		ao.ScheduleToCloseTimeout = time.Duration(*a.ScheduleToCloseTimeoutSeconds) * time.Second
	} else {
		ao.ScheduleToCloseTimeout = 24 * time.Hour
	}

	if a.ScheduleToStartTimeoutSeconds != nil {
		ao.ScheduleToStartTimeout = time.Duration(*a.ScheduleToStartTimeoutSeconds) * time.Second
	} else {
		ao.ScheduleToStartTimeout = 24 * time.Hour
	}

	ao.RetryPolicy = &temporal.RetryPolicy{
		MaximumAttempts:    1,
		InitialInterval:    0,
		BackoffCoefficient: 2,
		MaximumInterval:    1 * time.Minute,
	}

	if a.MaxRetries != nil {
		ao.RetryPolicy.MaximumAttempts = *a.MaxRetries
	}
	if a.RetryInitialIntervalSeconds != nil {
		ao.RetryPolicy.InitialInterval = time.Duration(*a.RetryInitialIntervalSeconds) * time.Second
	}
	if a.RetryBackoffCoefficient != nil {
		ao.RetryPolicy.BackoffCoefficient = *a.RetryBackoffCoefficient
	}
	if a.RetryMaxIntervalSeconds != nil {
		ao.RetryPolicy.MaximumInterval = time.Duration(*a.RetryMaxIntervalSeconds) * time.Second
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	args, err := makeInput(a.With, bindings, a.Key, a.SaveOutput, a.Input)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to make input.", "Error", err)
		return err
	}

	var result []byte
	err = workflow.ExecuteActivity(ctx, "Executor", a.Uses, args).Get(ctx, &result)
	if err != nil {
		return err
	}

	var output map[string]any
	if err := json.Unmarshal(result, &output); err != nil {
		if a.OnError != nil {
			return a.OnError.execute(ctx, bindings)
		}
		return err
	}

	saveOutput(output, bindings, a.Key)
	if a.OnSuccess != nil {
		return a.OnSuccess.execute(ctx, bindings)
	}

	return nil
}

func (s *Sequence) execute(ctx workflow.Context, bindings map[string]any) error {
	for _, stmt := range s.Elements {
		if err := stmt.execute(ctx, bindings); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parallel) execute(ctx workflow.Context, bindings map[string]any) error {
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error

	for _, stmt := range p.Branches {
		f := executeAsync(stmt, childCtx, bindings)
		selector.AddFuture(f, func(f workflow.Future) {
			if err := f.Get(ctx, nil); err != nil {
				cancelHandler()
				activityErr = err
			}
		})
	}

	for i := 0; i < len(p.Branches); i++ {
		selector.Select(ctx)
		if activityErr != nil {
			return activityErr
		}
	}
	return nil
}

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]any) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}
