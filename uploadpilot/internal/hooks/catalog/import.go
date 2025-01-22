package catalog

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
)

func (c *hooksCatalogService) AddMetadataHookFunc(ctx context.Context, input *hooks.HookInput, continueOnError bool) *hooks.HookResponse {
	imp := input.Import
	tusdHook := input.TusdHook

	// Extract Metadata
	var metadata map[string]interface{}
	err := mapstructure.Decode(tusdHook.Upload.MetaData, &metadata)
	if err != nil {
		msg := fmt.Sprintf("Failed to extract metadata from upload request: %s", err.Error())
		return c.updateImportAndReturnErrorResponse(ctx, input, errors.New(msg), continueOnError)
	}
	imp.Metadata = metadata

	return c.updateImportAndReturnSuccessResponse(ctx, input, "Metadata added to the import successfully")
}
