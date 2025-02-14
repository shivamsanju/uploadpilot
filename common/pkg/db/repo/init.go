package repo

import "github.com/uploadpilot/uploadpilot/common/pkg/db"

type Repositories struct {
	UserRepo            *UserRepo
	WorkspaceRepo       *WorkspaceRepo
	WorkspaceConfigRepo *WorkspaceConfigRepo
	WorkspaceUserRepo   *WorkspaceUserRepo
	UploadRepo          *UploadRepo
	UploadLogsRepo      *UploadLogsRepo
	ProcessorRepo       *ProcessorRepo
}

func NewRepositories(db *db.DB) *Repositories {
	return &Repositories{
		UserRepo:            NewUserRepo(db),
		WorkspaceRepo:       NewWorkspaceRepo(db),
		WorkspaceConfigRepo: NewWorkspaceConfigRepo(db),
		WorkspaceUserRepo:   NewWorkspaceUserRepo(db),
		UploadRepo:          NewUploadRepo(db),
		UploadLogsRepo:      NewUploadLogsRepo(db),
		ProcessorRepo:       NewProcessorRepo(db),
	}
}
