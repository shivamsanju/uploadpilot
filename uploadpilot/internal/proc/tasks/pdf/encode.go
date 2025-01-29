package pdf

import (
	"context"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type encodePDFContentTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewEncodePDFContentTask() tasks.Task {
	return &encodePDFContentTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *encodePDFContentTask) Do(ctx context.Context) error {
	t.leb.Publish(events.NewLogEvent(ctx, t.WorkspaceID, t.UploadID, "encoding pdf content", models.UploadLogLevelInfo))
	return nil
}
