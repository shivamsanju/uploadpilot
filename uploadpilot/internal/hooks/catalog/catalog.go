package catalog

import (
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
)

var (
	AuthenticateRequest = "AUTH_REQUEST"
	ValidateFile        = "VALIDATE_FILE"
	AddMetadata         = "ADD_METADATA"
	AddUploadURL        = "ADD_UPLOAD_URL"
	TriggerWebhook      = "TRIGGER_WEBHOOK"
)

type hooksCatalogService struct {
	WorkspaceRepo db.WorkspaceRepo
	ImportRepo    db.ImportRepo
	WebhookRepo   db.WebhookRepo
	Catalog       map[string]hooks.HookFn
}

func NewHooksCatalogService() *hooksCatalogService {
	svc := &hooksCatalogService{
		WorkspaceRepo: db.NewWorkspaceRepo(),
		ImportRepo:    db.NewImportRepo(),
		WebhookRepo:   db.NewWebhookRepo(),
		Catalog:       map[string]hooks.HookFn{},
	}

	svc.Catalog[AuthenticateRequest] = svc.AuthHookFunc
	svc.Catalog[AddMetadata] = svc.AddMetadataHookFunc
	svc.Catalog[AddUploadURL] = svc.AddUploadURLHookFunc
	svc.Catalog[TriggerWebhook] = svc.TriggerWebhookHookFunc

	return svc
}

func BuildPrefinishResponseHookExecutor() hooks.Executor {
	catalogSvc := NewHooksCatalogService()
	executor := hooks.NewHooksExecutor()

	executor.AddHook(&hooks.Hook{Name: AddMetadata, Execute: catalogSvc.AddMetadataHookFunc})
	executor.AddHook(&hooks.Hook{Name: AddUploadURL, Execute: catalogSvc.AddUploadURLHookFunc})

	return executor
}

func BuildPostfinishResponseHookExecutor() hooks.Executor {
	catalogSvc := NewHooksCatalogService()
	executor := hooks.NewHooksExecutor()

	executor.AddHook(&hooks.Hook{Name: TriggerWebhook, Execute: catalogSvc.TriggerWebhookHookFunc})

	return executor
}
