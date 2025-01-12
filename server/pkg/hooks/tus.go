package hooks

import (
	"fmt"
	"os"
	"path"

	"github.com/shivamsanju/uploader/pkg/globals"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func RemoveTusUploadHook(hook tusd.HookEvent) error {
	globals.Log.Infof("upload metadata -> %+v", hook.Upload.MetaData)
	uploadID := hook.Upload.ID
	err := os.Remove(path.Join(globals.TusUploadDir, fmt.Sprintf("%s.info", uploadID)))
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(globals.TusUploadDir, uploadID))
	if err != nil {
		return err
	}
	return nil
}
