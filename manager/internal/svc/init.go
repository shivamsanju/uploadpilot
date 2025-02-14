package svc

import "github.com/uploadpilot/uploadpilot/common/pkg/db/repo"

type Services struct {
	WorkspaceService *WorkspaceService
	UploadService    *UploadService
	UserService      *UserService
	ProcessorService *ProcessorService
}

func NewServices(repos *repo.Repositories) *Services {
	return &Services{
		WorkspaceService: NewWorkspaceService(repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.WorkspaceUserRepo, repos.UserRepo),
		UploadService:    NewUploadService(repos.UploadRepo, repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.UserRepo, repos.UploadLogsRepo),
		UserService:      NewUserService(repos.UserRepo),
		ProcessorService: NewProcessorService(repos.ProcessorRepo),
	}
}
