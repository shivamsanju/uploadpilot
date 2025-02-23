package img

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/phuslu/log"
	"github.com/uploadpilot/momentum/internal/data"
	"github.com/uploadpilot/momentum/internal/infra"
	"github.com/uploadpilot/momentum/internal/msg"
)

func BlurImage(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal activity arguments")
		return "", errors.New(msg.ErrInvalidActivityArguments)
	}

	radius, ok := args["radius"].(float64)
	if !ok {
		return "", errors.New("radius is required and must be a number")
	}

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
		img, err := imgio.Open(path)
		if err != nil {
			return err
		}
		blurredImg := blur.Gaussian(img, radius)
		return SaveImageBasedOnExtension(path, dh.OutputDir, blurredImg)
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}
