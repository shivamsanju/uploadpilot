package commontasks

import (
	"context"
	"fmt"

	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/infra"
	"github.com/uploadpilot/uploadpilot/manager/internal/msg"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

type ProcessingFailedTask struct {
	*tasks.BaseTask
	ueb    *events.UploadEventBus
	leb    *events.LogEventBus
	upRepo *db.UploadRepo
}

func NewFailedTask() tasks.Task {
	return &ProcessingFailedTask{
		ueb:      events.GetUploadEventBus(),
		leb:      events.GetLogEventBus(),
		upRepo:   db.NewUploadRepo(),
		BaseTask: &tasks.BaseTask{},
	}
}

func (t *ProcessingFailedTask) Do(ctx context.Context) error {
	infra.Log.Info("processing failed task...")

	uID := t.UploadID
	wID := t.WorkspaceID
	pID := t.ProcessorID
	tID := "WorkflowFailed"

	upload, err := t.upRepo.Get(ctx, uID)
	if err != nil {
		m := fmt.Sprintf("failed to mark processing status of upload [%s] as failed", uID)
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, &tID, models.UploadLogLevelError))
		return err
	}

	t.ueb.Publish(events.NewUploadEvent(ctx, events.EventUploadProcessingFailed, upload, "", nil))

	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, fmt.Sprintf(msg.ProcessingFailed, pID), &pID, &tID, models.UploadLogLevelInfo))
	return nil
}
