package services

import (
	"github.com/uploadpilot/core/internal/clients"
	"github.com/uploadpilot/core/internal/db/repo"
	"github.com/uploadpilot/core/internal/rbac"
)

type Services struct {
	TenantService    *TenantService
	WorkspaceService *WorkspaceService
	UploadService    *UploadService
	ProcessorService *ProcessorService
	APIKeyService    *APIKeyService
}

func NewServices(repos *repo.Repositories, clients *clients.Clients, accessManager *rbac.AccessManager) *Services {
	tenantSvc := NewTenantService(accessManager, repos.TenantRepo)
	workspaceSvc := NewWorkspaceService(accessManager, repos.WorkspaceRepo, repos.WorkspaceConfigRepo, clients.S3Client)
	apiKeySvc := NewAPIKeyService(accessManager, repos.APIKeyRepo, clients.KMSClient)
	processorSvc := NewProcessorService(accessManager, repos.ProcessorRepo, clients.TemporalClient, clients.S3Client)
	uploadSvc := NewUploadService(accessManager, repos.UploadRepo, workspaceSvc, processorSvc, clients.S3Client)

	return &Services{
		TenantService:    tenantSvc,
		WorkspaceService: workspaceSvc,
		UploadService:    uploadSvc,
		ProcessorService: processorSvc,
		APIKeyService:    apiKeySvc,
	}
}
