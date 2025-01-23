package upload

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/messages"
)

func (us *UploadService) extractMetadataFromTusdEvent(hook *tusd.HookEvent) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
	if err != nil {
		infra.Log.Errorf("Failed to extract metadata from upload request: %s", err.Error())
		return metadata, err
	}
	return metadata, nil
}

func (us *UploadService) getWorkspaceIDFromTusdEvent(hook *tusd.HookEvent) (string, error) {
	infra.Log.Infof("Upload Object %+v", hook.Upload)

	headers := hook.HTTPRequest.Header
	workspaceID := headers.Get("workspaceId")
	if len(workspaceID) == 0 {
		return "", fmt.Errorf(messages.InvalidWorkspaceIDInHeaders, workspaceID)
	}

	return workspaceID, nil
}
