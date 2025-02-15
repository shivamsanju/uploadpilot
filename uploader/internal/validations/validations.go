package validations

import (
	"fmt"
	"net/http"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/uploader/internal/msg"
)

func ValidateUploadSizeLimits(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) (string, error) {
	if config.MaxFileSize != nil && *config.MaxFileSize > 0 && hook.Upload.Size > int64(*config.MaxFileSize) {
		return "", fmt.Errorf(msg.MaxFileSizeExceeded, hook.Upload.Size, config.MaxFileSize)
	}

	if config.MinFileSize != nil && *config.MinFileSize > 0 && hook.Upload.Size < int64(*config.MinFileSize) {
		return "", fmt.Errorf(msg.MinFileSizeNotMet, hook.Upload.Size, config.MinFileSize)
	}

	maxFileSize := fmt.Sprintf("%d", config.MaxFileSize)
	if config.MaxFileSize != nil && *config.MaxFileSize == 0 {
		maxFileSize = "infinity"
	}

	return fmt.Sprintf(msg.UploadWithinSizeLimits, hook.Upload.Size, config.MinFileSize, maxFileSize), nil
}

func ValidateUploadFileType(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) (string, error) {
	if len(config.AllowedFileTypes) == 0 {
		return msg.FileTypeValidationSkipped, nil
	}

	for _, fileType := range config.AllowedFileTypes {
		if hook.Upload.MetaData["type"] == fileType {
			return fmt.Sprintf(msg.FileTypeValidationPassed, hook.Upload.MetaData["filetype"]), nil
		}
	}

	return "", fmt.Errorf(msg.FileTypeValidationFailed, hook.Upload.MetaData["type"])
}

func AuthenticateUpload(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) (string, error) {
	if config.AuthEndpoint == nil || *config.AuthEndpoint == "" {
		return msg.UploadAuthenticationSkipped, nil
	}

	req, err := http.NewRequestWithContext(hook.Context, "GET", *config.AuthEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf(msg.UploadAuthenticationFailed, err.Error())
	}

	tusdHeaders := hook.HTTPRequest.Header
	for key, values := range tusdHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf(msg.UploadAuthenticationFailed, err.Error())

	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf(msg.UploadAuthenticationFailed, "Error status code received from auth endpoint: "+resp.Status)
	}

	return msg.UploadAuthenticationPassed, nil
}
