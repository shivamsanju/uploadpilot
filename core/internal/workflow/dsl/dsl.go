package dsl

import (
	"encoding/json"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	WorkflowCtxKey string

	Workflow struct {
		Variables   map[string]any `json:"variables" yaml:"variables"`
		Root        Statement      `json:"root" yaml:"root"`
		WorkspaceID string         `json:"workspaceId"`
		UploadID    string         `json:"uploadId"`
		ProcessorID string         `json:"processorId"`
		FileName    string         `json:"fileName"`
		ContentType string         `json:"contentType"`
	}

	Statement struct {
		Activity  *ActivityInvocation `json:"activity,omitempty" yaml:"activity,omitempty"`
		Sequence  *Sequence           `json:"sequence,omitempty" yaml:"sequence,omitempty"`
		Parallel  *Parallel           `json:"parallel,omitempty" yaml:"parallel,omitempty"`
		Condition *Condition          `json:"condition,omitempty" yaml:"condition,omitempty"`
		Loop      *Loop               `json:"loop,omitempty" yaml:"loop,omitempty"`
	}

	Sequence struct {
		Elements []*Statement `json:"elements" yaml:"elements"`
	}

	Parallel struct {
		Branches []*Statement `json:"branches" yaml:"branches"`
	}

	Condition struct {
		Variable string     `json:"variable" yaml:"variable"`
		Value    string     `json:"value" yaml:"value"`
		Then     *Statement `json:"then" yaml:"then"`
		Else     *Statement `json:"else,omitempty" yaml:"else,omitempty"`
	}

	Loop struct {
		Iterations    int        `json:"iterations,omitempty" yaml:"iterations,omitempty"`
		Body          *Statement `json:"body" yaml:"body"`
		BreakVariable *string    `json:"break_variable,omitempty" yaml:"break_variable,omitempty"`
		BreakValue    *string    `json:"break_value,omitempty" yaml:"break_value,omitempty"`
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
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]any) error
	}
)

func (w Workflow) MarshalJSON() ([]byte, error) {
	type Alias Workflow
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&w),
	})
}

func (w Workflow) MarshalYAML() (interface{}, error) {
	type Alias Workflow
	return &struct {
		*Alias
	}{
		Alias: (*Alias)(&w),
	}, nil
}

func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	logger := workflow.GetLogger(ctx)
	bindings := make(map[string]any)
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	bindings["workspace_id"] = dslWorkflow.WorkspaceID
	bindings["upload_id"] = dslWorkflow.UploadID
	bindings["processor_id"] = dslWorkflow.ProcessorID
	bindings["file_name"] = dslWorkflow.FileName
	bindings["content_type"] = dslWorkflow.ContentType
	bindings["workflow_id"] = workflow.GetInfo(ctx).WorkflowExecution.ID
	bindings["run_id"] = workflow.GetInfo(ctx).WorkflowExecution.RunID

	defer func() {
		if err := doPostProcessing(ctx, bindings); err != nil {
			logger.Error("Error in post processing.", "Error", err)
		}
	}()

	err := dslWorkflow.Root.execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
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
	if b.Condition != nil {
		return b.Condition.execute(ctx, bindings)
	}
	if b.Loop != nil {
		return b.Loop.execute(ctx, bindings)
	}
	return nil
}

func (c *Condition) execute(ctx workflow.Context, bindings map[string]any) error {
	if bindings[c.Variable] == c.Value {
		if c.Then != nil {
			return c.Then.execute(ctx, bindings)
		}
	} else {
		if c.Else != nil {
			return c.Else.execute(ctx, bindings)
		}
	}
	return nil
}

func (l *Loop) execute(ctx workflow.Context, bindings map[string]any) error {
	for i := 0; i < l.Iterations; i++ {
		if err := l.Body.execute(ctx, bindings); err != nil {
			return err
		}
		if bindings[*l.BreakVariable] == *l.BreakValue {
			break
		}
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
		return err
	}

	saveOutput(output, bindings, a.Key)
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

func doPostProcessing(ctx workflow.Context, bindings map[string]any) error {
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
