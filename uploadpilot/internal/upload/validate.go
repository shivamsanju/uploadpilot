package upload

import (
	"errors"
	"fmt"
	"net/http"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/msg"
)

func (us *UploadService) ValidateUploadSizeLimits(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) error {

	if config.MaxFileSize > 0 && hook.Upload.Size > int64(config.MaxFileSize) {
		return fmt.Errorf(msg.MaxFileSizeExceeded, hook.Upload.Size, config.MaxFileSize)
	}

	if config.MinFileSize > 0 && hook.Upload.Size < int64(config.MinFileSize) {
		return fmt.Errorf(msg.MinFileSizeNotMet, hook.Upload.Size, config.MinFileSize)
	}

	maxFileSize := fmt.Sprintf("%d", config.MaxFileSize)
	if config.MaxFileSize == 0 {
		maxFileSize = "no_limit"
	}
	mesg := fmt.Sprintf("upload size: %d is within allowed range: %d to %s", hook.Upload.Size, config.MinFileSize, maxFileSize)
	us.logEventBus.Publish(events.NewLogEvent(hook.Context, workspaceID, uploadID, mesg, nil, nil, models.UploadLogLevelInfo))
	return nil
}

func (us *UploadService) ValidateUploadFileType(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) error {
	if len(config.AllowedFileTypes) == 0 {
		return nil
	}

	for _, fileType := range config.AllowedFileTypes {
		if hook.Upload.MetaData["type"] == fileType {
			mesg := fmt.Sprintf("file type: %s is among the allowed file types: %s", hook.Upload.MetaData["filetype"], config.AllowedFileTypes)
			us.logEventBus.Publish(events.NewLogEvent(hook.Context, workspaceID, uploadID, mesg, nil, nil, models.UploadLogLevelInfo))
			return nil
		}
	}

	return fmt.Errorf(msg.FileTypeNotAllowed, hook.Upload.MetaData["type"])
}

func (us *UploadService) AuthenticateUpload(hook *tusd.HookEvent, workspaceID, uploadID string, config *models.UploaderConfig) error {
	if len(config.AuthEndpoint) == 0 {
		return nil
	}

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

	mesg := fmt.Sprintf("upload is authenticated by auth endpoint: %s", config.AuthEndpoint)
	us.logEventBus.Publish(events.NewLogEvent(hook.Context, workspaceID, uploadID, mesg, nil, nil, models.UploadLogLevelInfo))
	return nil
}
