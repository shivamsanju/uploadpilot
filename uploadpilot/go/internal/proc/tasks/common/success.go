package commontasks

import (
	"context"
	"fmt"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type ProcessingSuccessTask struct {
	*tasks.BaseTask
	ueb    *events.UploadEventBus
	leb    *events.LogEventBus
	upRepo *db.UploadRepo
}

func NewSuccessTask() tasks.Task {
	return &ProcessingSuccessTask{
		ueb:      events.GetUploadEventBus(),
		leb:      events.GetLogEventBus(),
		upRepo:   db.NewUploadRepo(),
		BaseTask: &tasks.BaseTask{},
	}
}

func (t *ProcessingSuccessTask) Do(ctx context.Context) error {
	infra.Log.Info("processing success task...")

	uID := t.UploadID
	wID := t.WorkspaceID
	pID := t.ProcessorID
	tID := "WorkflowSuccess"

	upload, err := t.upRepo.Get(ctx, uID)
	if err != nil {
		m := fmt.Sprintf("failed to mark processing status of upload [%s] as success", uID)
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, &tID, models.UploadLogLevelError))
		return err
	}

	t.ueb.Publish(events.NewUploadEvent(ctx, events.EventUploadProcessed, upload, "", nil))

	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, fmt.Sprintf(msg.ProcessingComplete, pID), &pID, &tID, models.UploadLogLevelInfo))
	return nil
}
