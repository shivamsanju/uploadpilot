package zip

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/saracen/fastzip"
)

func ZipDirectory(ctx context.Context, sourceDir string, w io.Writer) error {
	a, err := fastzip.NewArchiver(w, sourceDir)
	if err != nil {
		return err
	}
	defer a.Close()

	files := make(map[string]os.FileInfo)
	if err = filepath.Walk(sourceDir, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	}); err != nil {
		return err
	}

	if err = a.Archive(ctx, files); err != nil {
		return err
	}

	return nil
}

func ZipFile(ctx context.Context, filePath string, w io.Writer) error {
	a, err := fastzip.NewArchiver(w, filePath)
	if err != nil {
		return err
	}
	defer a.Close()

	if err = a.Archive(ctx, nil); err != nil {
		return err
	}

	return nil
}

func UnzipFile(ctx context.Context, filePath, outDir string) error {
	e, err := fastzip.NewExtractor(filePath, outDir)
	if err != nil {
		return err
	}
	defer e.Close()

	if err = e.Extract(context.Background()); err != nil {
		return err
	}

	return nil
}
