package dsl

import (
	"encoding/json"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	WorkflowCtxKey string

	Workflow struct {
		Variables   map[string]interface{} `json:"variables" yaml:"variables"`
		Root        Statement              `json:"root" yaml:"root"`
		WorkspaceID string                 `json:"workspaceID" yaml:"workspaceID"`
		UploadID    string                 `json:"uploadID" yaml:"uploadID"`
		ProcessorID string                 `json:"processorID" yaml:"processorID"`
		FileName    string                 `json:"fileName" yaml:"fileName"`
		ContentType string                 `json:"contentType" yaml:"contentType"`
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
		BreakVariable *string    `json:"breakVariable,omitempty" yaml:"breakVariable,omitempty"`
		BreakValue    *string    `json:"breakValue,omitempty" yaml:"breakValue,omitempty"`
	}

	ActivityInvocation struct {
		Key                           string                 `json:"key" yaml:"key"`
		Uses                          string                 `json:"uses" yaml:"uses"`
		With                          map[string]interface{} `json:"with,omitempty" yaml:"with,omitempty"`
		Input                         *string                `json:"input,omitempty" yaml:"input,omitempty"`
		ScheduleToCloseTimeoutSeconds *int64                 `json:"scheduleToCloseTimeoutSeconds,omitempty" yaml:"scheduleToCloseTimeoutSeconds,omitempty"`
		ScheduleToStartTimeoutSeconds *int64                 `json:"scheduleToStartTimeoutSeconds,omitempty" yaml:"scheduleToStartTimeoutSeconds,omitempty"`
		StartToCloseTimeoutSeconds    *int64                 `json:"startToCloseTimeoutSeconds,omitempty" yaml:"startToCloseTimeoutSeconds,omitempty"`
		MaxRetries                    *int32                 `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`
		RetryBackoffCoefficient       *float64               `json:"retryBackoffCoefficient,omitempty" yaml:"retryBackoffCoefficient,omitempty"`
		RetryMaxIntervalSeconds       *int64                 `json:"retryMaxIntervalSeconds,omitempty" yaml:"retryMaxIntervalSeconds,omitempty"`
		RetryInitialIntervalSeconds   *int64                 `json:"retryInitialIntervalSeconds,omitempty" yaml:"retryInitialIntervalSeconds,omitempty"`
	}

	executable interface {
		execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error
	}
)

type WorkflowMeta struct {
	WorkspaceID string `json:"workspaceID" yaml:"workspaceID"`
	UploadID    string `json:"uploadID" yaml:"uploadID"`
	ProcessorID string `json:"processorID" yaml:"processorID"`
	WorkflowID  string `json:"workflowID" yaml:"workflowID"`
	RunID       string `json:"runID" yaml:"runID"`
	FileName    string `json:"fileName" yaml:"fileName"`
	ContentType string `json:"contentType" yaml:"contentType"`
}

func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	logger := workflow.GetLogger(ctx)
	bindings := make(map[string]interface{})
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	wfMeta := &WorkflowMeta{
		WorkspaceID: dslWorkflow.WorkspaceID,
		UploadID:    dslWorkflow.UploadID,
		ProcessorID: dslWorkflow.ProcessorID,
		FileName:    dslWorkflow.FileName,
		ContentType: dslWorkflow.ContentType,
		WorkflowID:  workflow.GetInfo(ctx).WorkflowExecution.ID,
		RunID:       workflow.GetInfo(ctx).WorkflowExecution.RunID,
	}

	wfMetaBytes, err := json.Marshal(wfMeta)
	if err != nil {
		return nil, err
	}
	wfMetaString := string(wfMetaBytes)

	err = dslWorkflow.Root.execute(ctx, wfMetaString, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}

func (b *Statement) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
	if b.Parallel != nil {
		return b.Parallel.execute(ctx, wfMetaString, bindings)
	}
	if b.Sequence != nil {
		return b.Sequence.execute(ctx, wfMetaString, bindings)
	}
	if b.Activity != nil {
		return b.Activity.execute(ctx, wfMetaString, bindings)
	}
	if b.Condition != nil {
		return b.Condition.execute(ctx, wfMetaString, bindings)
	}
	if b.Loop != nil {
		return b.Loop.execute(ctx, wfMetaString, bindings)
	}
	return nil
}

func (c *Condition) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
	if bindings[c.Variable] == c.Value {
		if c.Then != nil {
			return c.Then.execute(ctx, wfMetaString, bindings)
		}
	} else {
		if c.Else != nil {
			return c.Else.execute(ctx, wfMetaString, bindings)
		}
	}
	return nil
}

func (l *Loop) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
	for i := 0; i < l.Iterations; i++ {
		if err := l.Body.execute(ctx, wfMetaString, bindings); err != nil {
			return err
		}
		if bindings[*l.BreakVariable] == *l.BreakValue {
			break
		}
	}
	return nil
}

func (a *ActivityInvocation) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
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

	args, err := makeInput(a.With, bindings)
	if err != nil {
		logger := workflow.GetLogger(ctx)
		logger.Error("Failed to make input.", "Error", err)
		return err
	}

	var result []byte
	if a.Input == nil {
		a.Input = new(string)
	}

	err = workflow.ExecuteActivity(ctx, "Executor", a.Uses, args, wfMetaString).Get(ctx, &result)
	if err != nil {
		return err
	}
	bindings[a.Key+"result"] = result
	return nil
}

func (s *Sequence) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
	for _, stmt := range s.Elements {
		if err := stmt.execute(ctx, wfMetaString, bindings); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parallel) execute(ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) error {
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error

	for _, stmt := range p.Branches {
		f := executeAsync(stmt, childCtx, wfMetaString, bindings)
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

func executeAsync(exe executable, ctx workflow.Context, wfMetaString string, bindings map[string]interface{}) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, wfMetaString, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argMap map[string]interface{}, bindings map[string]interface{}) (string, error) {
	var args map[string]interface{}
	if argMap != nil {
		args = argMap
	} else {
		args = make(map[string]interface{})
	}
	for argument, value := range argMap {
		val, ok := value.(string)
		if ok && val[0] == '$' {
			varName := val[1:]
			args[argument] = bindings[varName]
		} else {
			args[argument] = value
		}
	}
	argsbytes, err := json.Marshal(args)
	if err != nil {
		return "", err
	}

	return string(argsbytes), nil
}
