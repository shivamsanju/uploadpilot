package img

import (
	"errors"
	"image"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
)

func SaveImageBasedOnExtension(path, outputDir string, processedImage image.Image) error {
	ext := filepath.Ext(path)
	filename := strings.TrimSuffix(filepath.Base(path), ext)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return imgio.Save(filepath.Join(outputDir, filename+".jpeg"), processedImage, imgio.JPEGEncoder(90))
	case ".png":
		return imgio.Save(filepath.Join(outputDir, filename+".png"), processedImage, imgio.PNGEncoder())
	case ".bmp":
		return imgio.Save(filepath.Join(outputDir, filename+".bmp"), processedImage, imgio.BMPEncoder())
	}
	return errors.New("unsupported image format")
}
