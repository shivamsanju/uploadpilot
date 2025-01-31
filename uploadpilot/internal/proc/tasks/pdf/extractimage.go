package pdf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type extractPDFImageTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewExtractPDFImageTask() tasks.Task {
	return &extractPDFImageTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *extractPDFImageTask) Do(ctx context.Context) error {
	t.Setup()
	defer t.Cleanup()

	wID := t.WorkspaceID
	uID := t.UploadID
	pID := t.ProcessorID
	tID := t.TaskID
	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, "extracting images from pdf", &pID, &tID, models.UploadLogLevelInfo))

	if err := t.SaveInputFile(ctx); err != nil {
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), &pID, &tID, models.UploadLogLevelError))
		return err
	}

	if err := t.extractPDFImage(); err != nil {
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), &pID, &tID, models.UploadLogLevelError))
		return err
	}
	objectName, err := t.SaveOutputFile(ctx)
	if err != nil {
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), &pID, &tID, models.UploadLogLevelError))
		return err
	}

	t.Output = map[string]interface{}{
		"inputObjId": objectName,
	}

	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, fmt.Sprintf("pdf image extracted to %s", objectName), &pID, &tID, models.UploadLogLevelInfo))

	return nil
}

func (t *extractPDFImageTask) extractPDFImage() error {
	inputDir := t.GetTaskInputDir()
	outDir := t.GetTaskOutDir()

	err := filepath.Walk(inputDir, func(pathname string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err := api.ExtractImagesFile(pathname, outDir, nil, nil); err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		return err
	}

	return nil
}
