package img

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/phuslu/log"
	"github.com/uploadpilot/momentum/internal/data"
	"github.com/uploadpilot/momentum/internal/infra"
	"github.com/uploadpilot/momentum/internal/msg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func ApplyTextWatermark(ctx context.Context, wfMeta, activityKey, inputActivityKey, argsStr string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal activity arguments")
		return "", errors.New(msg.ErrInvalidActivityArguments)
	}

	text, ok := args["text"].(string)
	if !ok || text == "" {
		return "", errors.New("text is required and must be a string")
	}

	opacity, ok := args["opacity"].(float64)
	if !ok || opacity < 0 || opacity > 1 {
		return "", errors.New("opacity is required and must be a number between 0 and 1")
	}

	s, ok := args["size"].(float64)
	if !ok {
		return "", errors.New("size is required and must be a number")
	}
	size := int(s)

	// Load TrueType font from gofont package
	fontBytes := goregular.TTF
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return "", errors.New("failed to parse built-in font")
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return "", errors.New("failed to create font face")
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
		bounds := img.Bounds()
		overlay := image.NewRGBA(bounds)
		draw.Draw(overlay, bounds, img, bounds.Min, draw.Src)

		// Set text color with opacity
		col := color.RGBA{255, 255, 255, uint8(255 * opacity)}

		// Set text position (20% of image width from the bottom right with margin)
		margin := 10
		point := fixed.Point26_6{
			X: fixed.I(bounds.Dx() - len(text)*size - margin),
			Y: fixed.I(bounds.Dy() - margin),
		}

		// Draw text using TrueType font
		d := &font.Drawer{
			Dst:  overlay,
			Src:  image.NewUniform(col),
			Face: face,
			Dot:  point,
		}
		d.DrawString(text)

		return SaveImageBasedOnExtension(path, dh.OutputDir, overlay)
	}); err != nil {
		return "", err
	}

	return dh.SaveOutput(ctx)
}
