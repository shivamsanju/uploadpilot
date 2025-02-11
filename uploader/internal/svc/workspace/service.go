package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db"
)

type WorkspaceService struct {
	workspaceRepo *db.WorkspaceRepo
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: db.NewWorkspaceRepo(),
	}
}

func (us *WorkspaceService) VerifySubscription(ctx context.Context, workspaceID string) (bool, error) {
	active, err := us.workspaceRepo.IsSubscriptionActive(ctx, workspaceID)
	return active, err
}
