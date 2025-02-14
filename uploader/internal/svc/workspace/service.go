package workspace

import (
	"context"

	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
)

type Service struct {
	workspaceRepo *repo.WorkspaceRepo
}

func NewWorkspaceService(workspaceRepo *repo.WorkspaceRepo) *Service {
	return &Service{
		workspaceRepo: workspaceRepo,
	}
}

func (us *Service) VerifySubscription(ctx context.Context, workspaceID string) (bool, error) {
	active, err := us.workspaceRepo.IsSubscriptionActive(ctx, workspaceID)
	return active, err
}
