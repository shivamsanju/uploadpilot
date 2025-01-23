package upload

import (
	"errors"
	"fmt"
	"net/http"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/messages"
)

func (us *UploadService) ValidateUploadSizeLimits(hook *tusd.HookEvent, config *models.UploaderConfig) error {

	if config.MaxFileSize > 0 && hook.Upload.Size > int64(config.MaxFileSize) {
		return fmt.Errorf(messages.MaxFileSizeExceeded, hook.Upload.Size, config.MaxFileSize)
	}

	if config.MinFileSize > 0 && hook.Upload.Size < int64(config.MinFileSize) {
		return fmt.Errorf(messages.MinFileSizeNotMet, hook.Upload.Size, config.MinFileSize)
	}

	return nil
}

func (us *UploadService) ValidateUploadFileType(hook *tusd.HookEvent, config *models.UploaderConfig) error {
	if len(config.AllowedFileTypes) == 0 {
		return nil
	}

	for _, fileType := range config.AllowedFileTypes {
		if hook.Upload.MetaData["type"] == fileType {
			return nil
		}
	}

	return fmt.Errorf(messages.FileTypeNotAllowed, hook.Upload.MetaData["type"])
}

func (us *UploadService) AuthenticateUpload(hook *tusd.HookEvent, config *models.UploaderConfig) error {
	req, err := http.NewRequestWithContext(hook.Context, "GET", config.AuthEndpoint, nil)
	if err != nil {
		return err
	}

	// Copy all headers from the original request to the new request
	tusdHeaders := hook.HTTPRequest.Header
	for key, values := range tusdHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err

	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(resp.Status)
	}

	return nil
}
