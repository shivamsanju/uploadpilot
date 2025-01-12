package hooks

import (
	"fmt"

	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func ValidateUpload(hook tusd.HookEvent) error {
	headers := hook.HTTPRequest.Header
	if uploaderId := headers.Get("uploaderId"); len(uploaderId) == 0 {
		return fmt.Errorf("missing uploaderId in header")
	}
	return nil
}
