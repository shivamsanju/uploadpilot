package repo

import "github.com/uploadpilot/go-core/db/pkg/driver"

type Repositories struct {
	TenantRepo          *TenantRepo
	SubscriptionRepo    *SubscriptionRepo
	WorkspaceRepo       *WorkspaceRepo
	WorkspaceConfigRepo *WorkspaceConfigRepo
	UploadRepo          *UploadRepo
	ProcessorRepo       *ProcessorRepo
	APIKeyRepo          *APIKeyRepo
	SecretsRepo         *SecretRepo
}

func NewRepositories(driver *driver.Driver) *Repositories {
	return &Repositories{
		TenantRepo:          NewTenantRepo(driver),
		SubscriptionRepo:    NewSubscriptionRepo(driver),
		WorkspaceRepo:       NewWorkspaceRepo(driver),
		WorkspaceConfigRepo: NewWorkspaceConfigRepo(driver),
		UploadRepo:          NewUploadRepo(driver),
		ProcessorRepo:       NewProcessorRepo(driver),
		APIKeyRepo:          NewAPIKeyRepo(driver),
		SecretsRepo:         NewSecretRepo(driver),
	}
}
