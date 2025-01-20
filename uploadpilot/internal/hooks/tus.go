package hooks

import (
	"fmt"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/uploadpilot/uploadpilot/internal/db"
)

func ValidateUpload(hook tusd.HookEvent) error {
	headers := hook.HTTPRequest.Header
	wsID := headers.Get("workspaceId")
	if len(wsID) == 0 {
		return fmt.Errorf("missing workspaceId in header")
	}
	workspaceID, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		return fmt.Errorf("invalid workspaceId: %w", err)
	}
	wsRepo := db.NewWorkspaceRepo()
	exists, err := wsRepo.CheckWorkspaceExists(hook.Context, workspaceID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("workspace not found")
	}
	return nil
}
