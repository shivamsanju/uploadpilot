package hooks

import (
	"context"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
)

type Status int

const (
	NotStarted Status = 0
	Running    Status = 1
	Success    Status = 2
	Failure    Status = 4
)

type HookInput struct {
	TusdHook *tusd.HookEvent
	Import   *models.Import
}

type HookResponse struct {
	Status          Status
	Error           error
	HookInput       *HookInput
	ContinueOnError bool
}

type HookFn func(ctx context.Context, input *HookInput, continueOnError bool) *HookResponse

type Hook struct {
	Name    string
	Execute HookFn
}
