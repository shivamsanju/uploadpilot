package catalog

import (
	"context"
	"reflect"
	"time"

	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *hooksCatalogService) updateImportAndReturnErrorResponse(ctx context.Context, input *hooks.HookInput, err error, continueOnError bool) *hooks.HookResponse {
	infra.Log.Errorf("[%s] Hook failed. Reason: %s. ContinueOnError: %t", reflect.TypeOf(input.TusdHook).Name(), err.Error(), continueOnError)

	input.Import.Logs = append(input.Import.Logs, models.Log{
		Message:   err.Error(),
		TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
	})

	if !continueOnError {
		input.Import.Status = models.ImportStatusFailed
	}

	_, err = c.ImportRepo.Update(ctx, input.Import.ID, input.Import)
	if err != nil {
		infra.Log.Errorf("Failed to update import  with ID: %s. Error: %s", input.Import.ID, err.Error())
	}

	return &hooks.HookResponse{
		Status:          hooks.Failure,
		Error:           err,
		ContinueOnError: continueOnError,
		HookInput:       input,
	}
}

func (c *hooksCatalogService) updateImportAndReturnSuccessResponse(ctx context.Context, input *hooks.HookInput, message string) *hooks.HookResponse {
	input.Import.Logs = append(input.Import.Logs, models.Log{
		Message:   message,
		TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
	})

	_, err := c.ImportRepo.Update(ctx, input.Import.ID, input.Import)
	if err != nil {
		infra.Log.Errorf("Failed to update import  with ID: %s. Error: %s", input.Import.ID, err.Error())
	}

	return &hooks.HookResponse{
		Status:          hooks.Success,
		Error:           nil,
		ContinueOnError: false,
		HookInput:       input,
	}
}
