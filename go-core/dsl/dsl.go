package dsl

import (
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
		Then     *Statement `json:"then,omitempty" yaml:"then,omitempty"`
		Else     *Statement `json:"else,omitempty" yaml:"else,omitempty"`
	}

	Loop struct {
		Iterations    int        `json:"iterations" yaml:"iterations"`
		Body          *Statement `json:"body" yaml:"body"`
		BreakVariable string     `json:"breakVariable" yaml:"breakVariable"`
		BreakValue    string     `json:"breakValue" yaml:"breakValue"`
	}

	ActivityInvocation struct {
		Name                          string   `json:"name" yaml:"name"`
		Arguments                     []string `json:"arguments" yaml:"arguments"`
		Result                        string   `json:"result" yaml:"result"`
		ScheduleToCloseTimeoutSeconds int64    `json:"scheduleToCloseTimeoutSeconds" yaml:"scheduleToCloseTimeoutSeconds"`
		ScheduleToStartTimeoutSeconds int64    `json:"scheduleToStartTimeoutSeconds" yaml:"scheduleToStartTimeoutSeconds"`
		StartToCloseTimeoutSeconds    int64    `json:"startToCloseTimeoutSeconds" yaml:"startToCloseTimeoutSeconds"`
		MaxRetries                    int32    `json:"maxRetries" yaml:"maxRetries"`
		RetryBackoffCoefficient       float64  `json:"retryBackoffCoefficient" yaml:"retryBackoffCoefficient"`
		RetryMaxIntervalSeconds       int64    `json:"retryMaxIntervalSeconds" yaml:"retryMaxIntervalSeconds"`
		RetryInitialIntervalSeconds   int64    `json:"retryInitialIntervalSeconds" yaml:"retryInitialIntervalSeconds"`
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
		if bindings[l.BreakVariable] == l.BreakValue {
			break
		}
	}
	return nil
}

func (a *ActivityInvocation) execute(ctx workflow.Context, bindings map[string]string) error {
	inputParam := makeInput(a.Arguments, bindings)
	var result string

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout:    time.Duration(a.StartToCloseTimeoutSeconds) * time.Second,
		ScheduleToCloseTimeout: time.Duration(a.ScheduleToCloseTimeoutSeconds) * time.Second,
		ScheduleToStartTimeout: time.Duration(a.ScheduleToStartTimeoutSeconds) * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Duration(a.RetryInitialIntervalSeconds) * time.Second,
			BackoffCoefficient: a.RetryBackoffCoefficient,
			MaximumInterval:    time.Duration(a.RetryMaxIntervalSeconds) * time.Second,
			MaximumAttempts:    a.MaxRetries,
		},
	})
	err := workflow.ExecuteActivity(ctx, a.Name, inputParam).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != "" {
		bindings[a.Result] = result
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
