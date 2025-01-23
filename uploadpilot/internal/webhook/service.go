package webhook

import (
	"github.com/uploadpilot/uploadpilot/internal/db"
)

type WebhookService struct {
	webRepo *db.WebhookRepo
	wsRepo  *db.WorkspaceRepo
	upRepo  *db.UploadRepo
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		webRepo: db.NewWebhookRepo(),
		wsRepo:  db.NewWorkspaceRepo(),
		upRepo:  db.NewUploadRepo(),
	}
}
