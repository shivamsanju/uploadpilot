package hooks

import (
	"fmt"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/pkg/db/repo"
)

func ValidateUpload(hook tusd.HookEvent) error {
	headers := hook.HTTPRequest.Header
	uploaderId := headers.Get("uploaderId")
	if len(uploaderId) == 0 {
		return fmt.Errorf("missing uploaderId in header")
	}
	uploaderRepo := repo.NewUploaderRepo()
	uploader, err := uploaderRepo.GetDataStoreCreds(hook.Context, uploaderId)
	if err != nil {
		return err
	}
	if uploader == nil {
		return fmt.Errorf("uploader not found")
	}
	return nil
}
