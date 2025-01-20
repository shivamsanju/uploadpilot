package hooks

import (
	"fmt"
	"os"
	"path"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

func RenameTusUploadHook(hook tusd.HookEvent, filename string) {
	uploadID := hook.Upload.ID
	err := os.Remove(path.Join(config.TusUploadDir, fmt.Sprintf("%s.info", uploadID)))
	if err != nil {
		infra.Log.Infof("failed to remove tus upload info file -> %+v", hook.Upload.ID)
	}
	err = os.Rename(path.Join(config.TusUploadDir, uploadID), path.Join(config.TusUploadDir, fmt.Sprintf("%s_%s", uploadID, filename)))
	if err != nil {
		infra.Log.Infof("failed to remove tus upload -> %+v", hook.Upload.ID)
	}
}

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
