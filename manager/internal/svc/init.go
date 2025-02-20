package svc

import (
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc/processor"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc/upload"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc/user"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc/workspace"
)

type Services struct {
	WorkspaceService *workspace.Service
	UploadService    *upload.Service
	UserService      *user.Service
	ProcessorService *processor.Service
}

func NewServices(repos *repo.Repositories) *Services {
	return &Services{
		WorkspaceService: workspace.NewService(repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.WorkspaceUserRepo, repos.UserRepo),
		UploadService:    upload.NewService(repos.UploadRepo, repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.UserRepo, repos.UploadLogsRepo),
		UserService:      user.NewService(repos.UserRepo),
		ProcessorService: processor.NewService(repos.ProcessorRepo),
	}
}
