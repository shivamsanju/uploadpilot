package notify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
)

type webhookTask struct {
	*tasks.BaseTask
	uploadRepo *db.UploadRepo
	leb        *events.LogEventBus
}

type webhookInput struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

func NewWebhookTask() tasks.Task {
	return &webhookTask{
		uploadRepo: db.NewUploadRepo(),
		leb:        events.GetLogEventBus(),
		BaseTask:   &tasks.BaseTask{},
	}
}

func (t *webhookTask) Do(ctx context.Context) error {
	wID := t.WorkspaceID
	uID := t.UploadID
	pID := t.ProcessorID
	tID := t.TaskID
	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, "triggering webhook", &pID, &tID, models.UploadLogLevelInfo))

	infra.Log.Infof("TaskDatasss: %+v", t.TaskData)
	webhook := &webhookInput{
		URL:    t.TaskData["url"].(string),
		Secret: t.TaskData["secret"].(string),
	}

	inputObjId, ok := t.Input["inputObjId"].(string)
	if !ok {
		m := fmt.Sprintf("missing required input: inputObjId for task %s", tID)
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, &tID, models.UploadLogLevelError))
		return errors.New(m)
	}

	assetUrl, err := s3.NewPresignClient(infra.S3Client).PresignGetObject(
		ctx,
		&s3.GetObjectInput{Bucket: &config.S3BucketName, Key: &inputObjId},
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(7 * 24 * time.Hour)
		},
	)
	if err != nil {
		m := fmt.Sprintf("failed to presign object for task %s: %s", tID, err.Error())
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, &tID, models.UploadLogLevelError))
		return err
	}

	if err := doWebhookRequest(webhook, assetUrl.URL); err != nil {
		m := fmt.Sprintf("failed to trigger webhook for task %s: %s", tID, err.Error())
		t.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, &tID, models.UploadLogLevelError))
		return err
	}

	message := fmt.Sprintf("Delivered to webhook [%s] successfully", webhook.URL)
	infra.Log.Info(message)
	t.leb.Publish(events.NewLogEvent(ctx, wID, uID, message, &pID, &tID, models.UploadLogLevelInfo))

	return nil
}

func doWebhookRequest(webhook *webhookInput, assetUrl string) error {
	body := fmt.Sprintf(`{
		"file_url": "%s",
	}`, assetUrl)

	infra.Log.Infof("triggering body {%s}", body)
	client := &http.Client{}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Secret", webhook.Secret)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("Error status code received from webhook: %s", resp.Status)
		return errors.New(msg)
	}

	return nil
}
