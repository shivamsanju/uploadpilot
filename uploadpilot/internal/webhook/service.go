package webhook

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"go.mongodb.org/mongo-driver/bson"
)

type WebhookService struct {
	webRepo      *db.WebhookRepo
	wsRepo       *db.WorkspaceRepo
	upRepo       *db.UploadRepo
	logsEventBus *events.LogEventBus
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		webRepo:      db.NewWebhookRepo(),
		wsRepo:       db.NewWorkspaceRepo(),
		upRepo:       db.NewUploadRepo(),
		logsEventBus: events.GetLogEventBus(),
	}
}

func (ws *WebhookService) GetAllWebhooksInWorkspace(ctx context.Context, workspaceID string) ([]models.Webhook, error) {
	return ws.webRepo.GetAll(ctx, workspaceID)
}

func (ws *WebhookService) GetWebhook(ctx context.Context, workspaceID, webhookID string) (*models.Webhook, error) {
	return ws.webRepo.Get(ctx, workspaceID, webhookID)
}

func (ws *WebhookService) CreateWebhook(ctx context.Context, workspaceID string, webhook *models.Webhook) error {
	if err := ws.validateWebhookBody(webhook); err != nil {
		return err
	}

	return ws.webRepo.Create(ctx, workspaceID, webhook)
}

func (ws *WebhookService) UpdateWebhook(ctx context.Context, workspaceID, webhookID string, webhook *models.Webhook) error {
	if err := ws.validateWebhookBody(webhook); err != nil {
		return err
	}

	return ws.webRepo.Update(ctx, workspaceID, webhookID, webhook)
}

func (ws *WebhookService) DeleteWebhook(ctx context.Context, workspaceID, webhookID string) error {
	return ws.webRepo.Delete(ctx, workspaceID, webhookID)
}

func (ws *WebhookService) PatchWebhook(ctx context.Context, workspaceID, webhookID string, patch *dto.PatchWebhookRequest) error {
	if reflect.TypeOf(patch.Enabled).Kind() != reflect.Bool {
		return fmt.Errorf("invalid type for enabled: %T", patch.Enabled)
	}

	return ws.webRepo.Patch(ctx, workspaceID, webhookID, bson.M{
		"enabled": patch.Enabled,
	})
}

func (ws *WebhookService) validateWebhookBody(webhook *models.Webhook) error {
	if err := infra.Validate.Struct(webhook); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return fmt.Errorf(msg.ValidationErr, errors)
	}
	return nil
}
