package dsl

import (
	"log"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	Workflow struct {
		Variables map[string]string `json:"variables" yaml:"variables"`
		Root      Statement         `json:"root" yaml:"root"`
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
		Name                          string   `json:"name" yaml:"name"`
		Arguments                     []string `json:"arguments,omitempty" yaml:"arguments,omitempty"`
		Result                        *string  `json:"result,omitempty" yaml:"result,omitempty"`
		ScheduleToCloseTimeoutSeconds *int64   `json:"scheduleToCloseTimeoutSeconds,omitempty" yaml:"scheduleToCloseTimeoutSeconds,omitempty"`
		ScheduleToStartTimeoutSeconds *int64   `json:"scheduleToStartTimeoutSeconds,omitempty" yaml:"scheduleToStartTimeoutSeconds,omitempty"`
		StartToCloseTimeoutSeconds    *int64   `json:"startToCloseTimeoutSeconds,omitempty" yaml:"startToCloseTimeoutSeconds,omitempty"`
		MaxRetries                    *int32   `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`
		RetryBackoffCoefficient       *float64 `json:"retryBackoffCoefficient,omitempty" yaml:"retryBackoffCoefficient,omitempty"`
		RetryMaxIntervalSeconds       *int64   `json:"retryMaxIntervalSeconds,omitempty" yaml:"retryMaxIntervalSeconds,omitempty"`
		RetryInitialIntervalSeconds   *int64   `json:"retryInitialIntervalSeconds,omitempty" yaml:"retryInitialIntervalSeconds,omitempty"`
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]string) error
	}
)

func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]string)
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	logger := workflow.GetLogger(ctx)

	err := dslWorkflow.Root.execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}

func (b *Statement) execute(ctx workflow.Context, bindings map[string]string) error {
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

func (c *Condition) execute(ctx workflow.Context, bindings map[string]string) error {
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

func (l *Loop) execute(ctx workflow.Context, bindings map[string]string) error {
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

func (a *ActivityInvocation) execute(ctx workflow.Context, bindings map[string]string) error {
	inputParam := makeInput(a.Arguments, bindings)
	var result string

	log.Printf("\n\n\n Wflow:  %+v \n\n\n", a)
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

	log.Printf("\n\n\n Wflow filled:  %+v \n\n\n", ao)

	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, a.Name, inputParam).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != nil {
		bindings[*a.Result] = result
	}
	return nil
}

func (s *Sequence) execute(ctx workflow.Context, bindings map[string]string) error {
	for _, stmt := range s.Elements {
		if err := stmt.execute(ctx, bindings); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parallel) execute(ctx workflow.Context, bindings map[string]string) error {
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

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]string) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argNames []string, argsMap map[string]string) []string {
	var args []string
	for _, arg := range argNames {
		args = append(args, argsMap[arg])
	}
	return args
}
