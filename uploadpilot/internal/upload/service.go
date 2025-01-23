package upload

import (
	"github.com/uploadpilot/uploadpilot/internal/db"
)

type UploadService struct {
	upRepo *db.UploadRepo
	wsRepo *db.WorkspaceRepo
}

func NewUploadService() *UploadService {
	return &UploadService{
		upRepo: db.NewUploadRepo(),
		wsRepo: db.NewWorkspaceRepo(),
	}
}
