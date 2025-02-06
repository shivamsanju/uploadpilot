package image

import (
	"context"
	"errors"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type extractTextFromImageTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewExtractTextFromImageTask() tasks.Task {
	return &extractTextFromImageTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *extractTextFromImageTask) Do(ctx context.Context) error {
	t.leb.Publish(events.NewLogEvent(ctx, t.WorkspaceID, t.UploadID, "OCR task is in developmet.", &t.ProcessorID, &t.TaskID, models.UploadLogLevelError))
	return errors.New("ocr task is in developmet")
}
