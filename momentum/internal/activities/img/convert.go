package img

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/uploadpilot/momentum/internal/data"
	"github.com/uploadpilot/momentum/internal/infra"
)

func ConvertImageToPng(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
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
		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		return imgio.Save(filepath.Join(dh.OutputDir, filename+".png"), img, imgio.PNGEncoder())
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}

func ConvertImageToJpeg(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
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
		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		return imgio.Save(filepath.Join(dh.OutputDir, filename+".jpeg"), img, imgio.JPEGEncoder(100))
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}

func ConvertImageToBmp(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
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
		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		return imgio.Save(filepath.Join(dh.OutputDir, filename+".bmp"), img, imgio.BMPEncoder())
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}
