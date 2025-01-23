package webhook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type webhookRoutineResponse struct {
	URL   string
	Error error
}

func (c *WebhookService) TriggerWebhook(upload *models.Upload) error {
	ctx := context.Background()
	workspaceID := upload.WorkspaceID.Hex()

	webhooks, err := c.webRepo.GetEnabledWebhooksWithSecret(ctx, workspaceID)
	if err != nil {
		return err
	}

	resultChan := make(chan *webhookRoutineResponse, len(webhooks))
	wg := &sync.WaitGroup{}
	for i := 0; i < len(webhooks); i++ {
		wg.Add(1)
		go triggerWebhookRoutine(wg, &webhooks[i], upload, resultChan)
	}
	wg.Wait()

	logs := []models.Log{}
	for i := 0; i < len(webhooks); i++ {
		res := <-resultChan
		now := primitive.NewDateTimeFromTime(time.Now())
		if res.Error != nil {
			logs = append(logs, models.Log{
				Message:   fmt.Sprintf("Error while triggering webhook [%s]: %s", res.URL, res.Error.Error()),
				TimeStamp: now,
			})
		} else {
			logs = append(logs, models.Log{
				Message:   fmt.Sprintf("Delivered to webhook [%s] successfully", res.URL),
				TimeStamp: now,
			})
		}
	}

	if len(logs) > 0 {
		if err := c.upRepo.AddLogs(ctx, workspaceID, upload.ID.Hex(), logs); err != nil {
			return err
		}
	}
	return nil
}

func triggerWebhookRoutine(wg *sync.WaitGroup, webhook *models.Webhook, upload *models.Upload, resultChan chan *webhookRoutineResponse) {
	defer wg.Done()

	body := fmt.Sprintf(`{
		"file_url": "%s",
		"file_name": "%s",
		"file_size": %d,
		"upload_time": "%s"
	}`, upload.URL, upload.Metadata["filename"], upload.Size, upload.StartedAt.Time().Format(time.RFC3339))

	infra.Log.Infof("triggering body {%s}", body)
	client := &http.Client{}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: err}
		return
	}
	req.Header.Set("Secret", webhook.SigningSecret)

	resp, err := client.Do(req)
	if err != nil {
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: err}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		msg := fmt.Sprintf("Error status code received from webhook: %s", resp.Status)
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: errors.New(msg)}
		return
	}
	resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: nil}
}
