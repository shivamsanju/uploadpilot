package svc

import (
	"github.com/uploadpilot/go-core/db/pkg/repo"
	"github.com/uploadpilot/manager/internal/svc/processor"
	"github.com/uploadpilot/manager/internal/svc/upload"
	"github.com/uploadpilot/manager/internal/svc/user"
	"github.com/uploadpilot/manager/internal/svc/workspace"
)

type Services struct {
	WorkspaceService *workspace.Service
	UploadService    *upload.Service
	UserService      *user.Service
	ProcessorService *processor.Service
}

func NewServices(repos *repo.Repositories) *Services {
	processorSvc := processor.NewService(repos.ProcessorRepo)
	workspaceSvc := workspace.NewService(repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.WorkspaceUserRepo, repos.UserRepo)
	userSvc := user.NewService(repos.UserRepo)
	uploadSvc := upload.NewService(repos.UploadRepo, repos.WorkspaceRepo, repos.WorkspaceConfigRepo, repos.UserRepo, processorSvc)

	return &Services{
		WorkspaceService: workspaceSvc,
		UploadService:    uploadSvc,
		UserService:      userSvc,
		ProcessorService: processorSvc,
	}
}
