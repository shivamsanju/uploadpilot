package hooks

import (
	"context"
	"fmt"
)

// Executor
type Executor interface {
	Start(ctx context.Context, input *HookInput, continueOnError bool) error
	AddHook(hook *Hook)
}

type executor struct {
	hooks []*Hook
}

func NewHooksExecutor() Executor {
	return &executor{
		hooks: []*Hook{},
	}
}

func (e *executor) AddHook(hook *Hook) {
	e.hooks = append(e.hooks, hook)
}

func (e *executor) Start(ctx context.Context, input *HookInput, continueOnError bool) error {
	for _, hook := range e.hooks {
		hookResp := hook.Execute(ctx, input, continueOnError)
		if hookResp.Error != nil && !hookResp.ContinueOnError {
			return fmt.Errorf("hook [%s] failed: %w", hook.Name, hookResp.Error)
		}
	}

	return nil
}
