package catalog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type webhookRoutineResponse struct {
	URL   string
	Error error
}

func (c *hooksCatalogService) TriggerWebhookHookFunc(ctx context.Context, input *hooks.HookInput, continueOnError bool) *hooks.HookResponse {
	workspaceID := input.Import.WorkspaceID
	input.Import.Status = models.ImportStatusSuccess
	webhooks, err := c.WebhookRepo.GetEnabledWebhooksWithSecret(ctx, workspaceID)
	if err != nil {
		return c.updateImportAndReturnErrorResponse(ctx, input, err, continueOnError)
	}

	resultChan := make(chan *webhookRoutineResponse, len(webhooks))
	wg := &sync.WaitGroup{}
	for _, webhook := range webhooks {
		wg.Add(1)
		go triggerWebhookRoutine(wg, &webhook, input, resultChan)
	}
	wg.Wait()

	for i := 0; i < len(webhooks); i++ {
		res := <-resultChan
		if res.Error != nil {
			// TODO: add to logs
			// input.Import.Status = models.ImportStatusFailed
			input.Import.Logs = append(input.Import.Logs, models.Log{
				Message:   fmt.Sprintf("Error while triggering webhook [%s]: %s", res.URL, res.Error.Error()),
				TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
			})
		} else {
			input.Import.Logs = append(input.Import.Logs, models.Log{
				Message:   fmt.Sprintf("Delivered to webhook [%s] successfully", res.URL),
				TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
			})
		}
	}
	return c.updateImportAndReturnSuccessResponse(ctx, input, "Webhooks triggered successfully")
}

func triggerWebhookRoutine(wg *sync.WaitGroup, webhook *models.Webhook, input *hooks.HookInput, resultChan chan *webhookRoutineResponse) {
	defer wg.Done()
	metadata := map[string]string{}
	err := mapstructure.Decode(input.Import.Metadata, &metadata)
	if err != nil {
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: err}
		return
	}
	body := fmt.Sprintf(`{
		"file_url": "%s",
		"file_name": "%s",
		"file_size": %d,
		"upload_time": "%s"
	}`, input.Import.URL, metadata["filename"], input.Import.Size, input.Import.StartedAt.Time().Format(time.RFC3339))

	infra.Log.Infof("triggering body {%s}", body)
	client := &http.Client{}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		msg := fmt.Sprintf("Error in building http request for authenticating upload request: %s", err.Error())
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: errors.New(msg)}
		return
	}

	req.Header.Set("Secret", webhook.SigningSecret)

	aResp, err := client.Do(req)
	if err != nil {
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: err}
		return
	}
	defer aResp.Body.Close()
	infra.Log.Infof("status code received while triggering webhook {%s}: %d", webhook.URL, aResp.StatusCode)
	if aResp.StatusCode < 200 || aResp.StatusCode > 299 {
		msg := fmt.Sprintf("Erroneous status code received while triggering webhook: %s", aResp.Status)
		resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: errors.New(msg)}
		return
	}
	resultChan <- &webhookRoutineResponse{URL: webhook.URL, Error: nil}
}
