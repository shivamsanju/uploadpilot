package trigger

import (
	"context"

	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

type triggerTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	procRepo   *db.ProcessorRepo
	leb        *events.LogEventBus
}

func NewTriggerTask() tasks.Task {
	return &triggerTask{
		uploadRepo: db.NewUploadRepo(),
		procRepo:   db.NewProcessorRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *triggerTask) Do(ctx context.Context) error {
	t.Output = t.Input
	return nil
}
