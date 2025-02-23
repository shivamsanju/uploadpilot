package img

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanoberholster/imagemeta"
	"github.com/uploadpilot/momentum/internal/data"
	"github.com/uploadpilot/momentum/internal/infra"
)

func ExtractMetadataFromImage(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	dh := data.NewActivityDataHandler(ctx, wfMeta, activityKey, inputActivityKey, infra.S3Client)
	if err := dh.PrepareDataLayer(ctx); err != nil {
		return "", err
	}
	defer dh.Cleanup(ctx)

	if err := filepath.Walk(dh.InputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		meta, err := imagemeta.Decode(file)
		if err != nil {
			return err
		}

		metadata, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return err
		}

		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		return os.WriteFile(filepath.Join(dh.OutputDir, filename+".json"), metadata, 0644)
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}
