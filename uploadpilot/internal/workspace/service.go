package workspace

import (
	"github.com/uploadpilot/uploadpilot/internal/db"
)

type WorkspaceService struct {
	wsRepo   *db.WorkspaceRepo
	userRepo *db.UserRepo
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		wsRepo:   db.NewWorkspaceRepo(),
		userRepo: db.NewUserRepo(),
	}
}
