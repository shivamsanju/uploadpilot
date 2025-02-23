package img

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/phuslu/log"
	"github.com/uploadpilot/momentum/internal/data"
	"github.com/uploadpilot/momentum/internal/infra"
	"github.com/uploadpilot/momentum/internal/msg"
)

func ResizeImage(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal activity arguments")
		return "", errors.New(msg.ErrInvalidActivityArguments)
	}

	// print type of args
	width, ok := args["width"].(float64)
	if !ok {
		return "", errors.New("width is required and must be a integer")
	}

	height, ok := args["height"].(float64)
	if !ok {
		return "", errors.New("height is required and must be a integer")
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
		resizedImg := transform.Resize(img, int(width), int(height), transform.Linear)
		if err := SaveImageBasedOnExtension(path, dh.OutputDir, resizedImg); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}
