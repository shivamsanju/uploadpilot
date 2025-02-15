package svc

import (
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/repo"
	uploaderconfig "github.com/uploadpilot/uploadpilot/uploader/internal/svc/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/processor"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/upload"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/workspace"
)

type Services struct {
	UploadService    *upload.Service
	WorkspaceService *workspace.Service
	ProcessorService *processor.Service
	ConfigService    *uploaderconfig.Service
}

func NewServices(repos *repo.Repositories) *Services {
	return &Services{
		UploadService:    upload.NewUploadService(repos.UploadRepo, repos.UploadLogsRepo, repos.WorkspaceRepo, repos.WorkspaceConfigRepo),
		WorkspaceService: workspace.NewWorkspaceService(repos.WorkspaceRepo),
		ProcessorService: processor.NewProcessorService(repos.ProcessorRepo),
		ConfigService:    uploaderconfig.NewConfigService(repos.WorkspaceConfigRepo),
	}
}
