package commontasks

import (
	"context"
	"os"

	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

type ProcessingCleanupTask struct {
	*tasks.BaseTask
}

func NewCleanupTask() tasks.Task {
	return &ProcessingCleanupTask{
		BaseTask: &tasks.BaseTask{},
	}
}

func (t *ProcessingCleanupTask) Do(ctx context.Context) error {
	infra.Log.Info("processing cleanup task...")
	t.TmpDirLock.Lock()
	err := os.RemoveAll(t.TmpDir)
	if err != nil {
		infra.Log.Errorf("failed to remove tmp dir: %s", err.Error())
	}
	t.TmpDirLock.Unlock()
	return nil
}
