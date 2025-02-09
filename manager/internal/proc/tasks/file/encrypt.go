package file

import (
	"context"
	"errors"

	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

type encryptContentTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewEncryptContentTask() tasks.Task {
	return &encryptContentTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *encryptContentTask) Do(ctx context.Context) error {
	t.leb.Publish(events.NewLogEvent(ctx, t.WorkspaceID, t.UploadID, "encrypting content task is in developmet", &t.ProcessorID, &t.TaskID, models.UploadLogLevelInfo))
	return errors.New("encrypting content task is in developmet")
}
