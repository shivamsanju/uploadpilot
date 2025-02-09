package pdf

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ledongthuc/pdf"
	"github.com/uploadpilot/uploadpilot/manager/internal/db"
	"github.com/uploadpilot/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/events"
	"github.com/uploadpilot/uploadpilot/manager/internal/proc/tasks"
)

type extractPDFContentTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

func NewExtractPDFContentTask() tasks.Task {
	return &extractPDFContentTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *extractPDFContentTask) Do(ctx context.Context) error {
	t.Setup()
	defer t.Cleanup()

	wID := t.WorkspaceID
	uID := t.UploadID
	pID := t.ProcessorID
	tID := t.TaskID
	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, "extracting pdf content", &pID, &tID, models.UploadLogLevelInfo))

	if err := t.SaveInputFile(ctx); err != nil {
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), &pID, &tID, models.UploadLogLevelError))
		return err
	}

	if err := t.extractPDFContent(); err != nil {
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

	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, fmt.Sprintf("pdf content extracted to %s", objectName), &pID, &tID, models.UploadLogLevelInfo))

	return nil
}

func (t *extractPDFContentTask) extractPDFContent() error {
	inputDir := t.GetTaskInputDir()
	outDir := t.GetTaskOutDir()

	err := filepath.Walk(inputDir, func(pathname string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		txt, err := readPdf(pathname)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(outDir, filepath.Base(pathname)+".txt"), []byte(txt), os.ModePerm); err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

func readPdf(path string) (string, error) {
	_, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	var textBuilder bytes.Buffer
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		txt, err := p.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		textBuilder.WriteString(txt)
	}
	return textBuilder.String(), nil
}
