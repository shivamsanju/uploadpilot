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
		BaseTask:   tasks.NewBaseTask(),
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
	}
}

func (t *encodePDFContentTask) Do(ctx context.Context) error {
	t.leb.Publish(events.NewLogEvent(ctx, t.Data.WorkspaceID, t.Data.UploadID, "encoding pdf content", models.UploadLogLevelInfo))
	return nil
}
