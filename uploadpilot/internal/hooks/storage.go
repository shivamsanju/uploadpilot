package hooks

import (
	"fmt"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateImportMetadata(hook tusd.HookEvent, imp *models.Import) error {
	headers := hook.HTTPRequest.Header
	wsID := headers.Get("workspaceId")
	if len(wsID) == 0 {
		return fmt.Errorf("missing workspaceId in header")
	}
	workspaceID, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		return fmt.Errorf("invalid workspaceId: %w", err)
	}
	imp.WorkspaceID = workspaceID
	objectFileName := hook.Upload.ID + "_" + hook.Upload.MetaData["filename"]
	imp.StoredFileName = objectFileName

	return nil
}
