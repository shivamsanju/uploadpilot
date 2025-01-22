package catalog

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/uploadpilot/uploadpilot/internal/hooks"
)

func (c *hooksCatalogService) AuthHookFunc(ctx context.Context, input *hooks.HookInput, continueOnError bool) *hooks.HookResponse {
	tusdHook := input.TusdHook

	headers := tusdHook.HTTPRequest.Header
	token := headers.Get("Authorization")
	if len(token) == 0 {
		msg := fmt.Sprintf("missing Authorization header in upload request: %s", tusdHook.HTTPRequest.Header)
		return c.updateImportAndReturnErrorResponse(ctx, input, errors.New(msg), continueOnError)
	}

	// Get auth hook - send a http request
	req, err := http.NewRequest("GET", "http://localhost:3001/api/auth", nil)
	if err != nil {
		msg := fmt.Sprintf("Error in building http request for authenticating upload request: %s", err.Error())
		return c.updateImportAndReturnErrorResponse(ctx, input, errors.New(msg), continueOnError)
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{}
	aResp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Error from http request: %s", err.Error())
		return c.updateImportAndReturnErrorResponse(ctx, input, errors.New(msg), continueOnError)
	}
	defer aResp.Body.Close()
	if aResp.StatusCode < 200 || aResp.StatusCode > 299 {
		msg := fmt.Sprintf("Erroneous status code received while authenticating upload request: %s", aResp.Status)
		return c.updateImportAndReturnErrorResponse(ctx, input, errors.New(msg), continueOnError)
	}

	return c.updateImportAndReturnSuccessResponse(ctx, input, "Authenticated upload request successfully")
}
