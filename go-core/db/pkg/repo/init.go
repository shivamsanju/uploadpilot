package repo

import "github.com/uploadpilot/uploadpilot/go-core/db/pkg/driver"

type Repositories struct {
	UserRepo            *UserRepo
	WorkspaceRepo       *WorkspaceRepo
	WorkspaceConfigRepo *WorkspaceConfigRepo
	WorkspaceUserRepo   *WorkspaceUserRepo
	UploadRepo          *UploadRepo
	UploadLogsRepo      *UploadLogsRepo
	ProcessorRepo       *ProcessorRepo
}

func NewRepositories(driver *driver.Driver) *Repositories {
	return &Repositories{
		UserRepo:            NewUserRepo(driver),
		WorkspaceRepo:       NewWorkspaceRepo(driver),
		WorkspaceConfigRepo: NewWorkspaceConfigRepo(driver),
		WorkspaceUserRepo:   NewWorkspaceUserRepo(driver),
		UploadRepo:          NewUploadRepo(driver),
		UploadLogsRepo:      NewUploadLogsRepo(driver),
		ProcessorRepo:       NewProcessorRepo(driver),
	}
}
