package svc

import (
	"github.com/uploadpilot/go-core/common/vault"
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/svc/apikey"
	"github.com/uploadpilot/manager/internal/svc/processor"
	"github.com/uploadpilot/manager/internal/svc/tenant"
	"github.com/uploadpilot/manager/internal/svc/upload"
	"github.com/uploadpilot/manager/internal/svc/workspace"
)

type Services struct {
	TenantService    *tenant.Service
	WorkspaceService *workspace.Service
	UploadService    *upload.Service
	ProcessorService *processor.Service
	APIKeyService    *apikey.Service
}

func NewServices(repos *repo.Repositories, kms vault.KMS) *Services {
	tenantSvc := tenant.NewService(repos.TenantRepo)
	workspaceSvc := workspace.NewService(repos.WorkspaceRepo, repos.WorkspaceConfigRepo)
	apiKeySvc := apikey.NewService(repos.APIKeyRepo, kms)
	processorSvc := processor.NewService(repos.ProcessorRepo)
	uploadSvc := upload.NewService(repos.UploadRepo, processorSvc)

	return &Services{
		TenantService:    tenantSvc,
		WorkspaceService: workspaceSvc,
		UploadService:    uploadSvc,
		ProcessorService: processorSvc,
		APIKeyService:    apiKeySvc,
	}
}
