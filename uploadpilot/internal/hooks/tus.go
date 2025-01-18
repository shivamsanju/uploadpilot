package hooks

import (
	"fmt"
	"os"
	"path"

	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/infra"
)

func RemoveTusUploadHook(hook tusd.HookEvent) {
	uploadID := hook.Upload.ID
	err := os.Remove(path.Join(config.TusUploadDir, fmt.Sprintf("%s.info", uploadID)))
	if err != nil {
		infra.Log.Infof("failed to remove tus upload -> %+v", hook.Upload.ID)
	}
	err = os.Remove(path.Join(config.TusUploadDir, uploadID))
	if err != nil {
		infra.Log.Infof("failed to remove tus upload -> %+v", hook.Upload.ID)
	}
}

func ValidateUpload(hook tusd.HookEvent) error {
	headers := hook.HTTPRequest.Header
	uploaderId := headers.Get("uploaderId")
	if len(uploaderId) == 0 {
		return fmt.Errorf("missing uploaderId in header")
	}
	uploaderRepo := db.NewUploaderRepo()
	uploader, err := uploaderRepo.GetDataStoreCreds(hook.Context, uploaderId)
	if err != nil {
		return err
	}
	if uploader == nil {
		return fmt.Errorf("uploader not found")
	}
	return nil
}
