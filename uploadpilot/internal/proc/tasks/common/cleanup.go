package commontasks

import (
	"context"
	"os"

	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type ProcessingCleanupTask struct {
	*tasks.BaseTask
}

func NewCleanupTask() tasks.Task {
	return &ProcessingCleanupTask{
		BaseTask: tasks.NewBaseTask(),
	}
}

func (t *ProcessingCleanupTask) Do(ctx context.Context) error {
	infra.Log.Info("processing cleanup task...")
	t.Data.TmpDirLock.Lock()
	os.RemoveAll(t.Data.TmpDir)
	t.Data.TmpDirLock.Unlock()
	return nil
}
