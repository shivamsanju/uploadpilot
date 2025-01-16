package hooks

import (
	"fmt"
	"os"
	"path"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/pkg/globals"
)

func RemoveTusUploadHook(hook tusd.HookEvent) {
	uploadID := hook.Upload.ID
	err := os.Remove(path.Join(globals.TusUploadDir, fmt.Sprintf("%s.info", uploadID)))
	if err != nil {
		globals.Log.Infof("failed to remove tus upload -> %+v", hook.Upload.ID)
	}
	err = os.Remove(path.Join(globals.TusUploadDir, uploadID))
	if err != nil {
		globals.Log.Infof("failed to remove tus upload -> %+v", hook.Upload.ID)
	}
}
