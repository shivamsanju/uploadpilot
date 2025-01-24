package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/messages"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (us *UploadService) GetAllUploads(ctx context.Context, workspaceID string, skip int64, limit int64, search string) ([]models.Upload, int64, error) {
	if strings.HasPrefix(search, "{") {
		searchParams, err := utils.ExtractKeyValuePairs(search)
		if err != nil {
			return nil, 0, err
		}
		return us.upRepo.GetAllFilterByMetadata(ctx, workspaceID, skip, limit, searchParams)
	}

	return us.upRepo.GetAll(ctx, workspaceID, skip, limit, search)
}

func (us *UploadService) GetUploadDetails(ctx context.Context, workspaceID, uploadID string) (*models.Upload, error) {
	return us.upRepo.Get(ctx, workspaceID, uploadID)
}

func (us *UploadService) GetUploadDetailsFromTusdEvent(hook *tusd.HookEvent) (*models.Upload, error) {
	uploadID := hook.Upload.MetaData["uploadId"]
	hook.Upload.ID = uploadID
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return nil, err
	}
	return us.upRepo.Get(hook.Context, workspaceID, uploadID)
}

func (us *UploadService) CreateUpload(hook *tusd.HookEvent) (*models.Upload, error) {
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return nil, err
	}

	// workspace existence is anyways checked here, no need to add extra method call
	config, err := us.getUploaderConfig(hook)
	if err != nil {
		return nil, err
	}

	// beyond this point, an upload attempt will be logged
	now := primitive.NewDateTimeFromTime(time.Now())
	upload := &models.Upload{
		Size:      hook.Upload.Size,
		Status:    models.UploadStatusInProgress,
		StartedAt: now,
		Logs: []models.Log{{
			Message:   "File upload started",
			TimeStamp: now,
		}},
	}

	logErrorAndUpdateUpload := func(err error) error {
		now := primitive.NewDateTimeFromTime(time.Now())
		upload.FinishedAt = now
		upload.Status = models.UploadStatusFailed
		upload.Logs = append(upload.Logs, models.Log{Message: err.Error(), TimeStamp: now})
		if repoErr := us.upRepo.Create(hook.Context, workspaceID, upload); repoErr != nil {
			infra.Log.Errorf("unable to create upload: %s", repoErr)
			return repoErr
		}
		return err
	}

	metadata, err := us.extractMetadataFromTusdEvent(hook)
	if err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}
	upload.Metadata = metadata

	if err := us.ValidateUploadSizeLimits(hook, config); err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}
	upload.Logs = append(upload.Logs, models.Log{Message: "Upload size is within the allowed range", TimeStamp: now})

	if err := us.ValidateUploadFileType(hook, config); err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}
	upload.Logs = append(upload.Logs, models.Log{Message: "Uploaded file type is allowed", TimeStamp: now})

	if len(config.AuthEndpoint) > 0 {
		if err := us.AuthenticateUpload(hook, config); err != nil {
			return nil, logErrorAndUpdateUpload(fmt.Errorf(messages.UploadAuthenticationFailed, err))
		}
		upload.Logs = append(upload.Logs, models.Log{Message: "Authenticated upload request successfully", TimeStamp: now})
	}

	if err := us.upRepo.Create(hook.Context, workspaceID, upload); err != nil {
		infra.Log.Errorf("unable to create upload: %s", err)
		return nil, err
	}

	return upload, nil
}

func (us *UploadService) FinishUpload(ctx context.Context, upload *models.Upload) error {
	now := primitive.NewDateTimeFromTime(time.Now())
	url, storedFileName, err := us.GenerateS3URLFromUploadID(ctx, upload.ID.Hex())
	if err != nil {
		upload.Status = models.UploadStatusFailed
		upload.Logs = append(upload.Logs, models.Log{Message: err.Error(), TimeStamp: now})
	} else {
		upload.URL = url
		upload.Status = models.UploadStatusSuccess
		upload.FinishedAt = now
		upload.StoredFileName = storedFileName
		upload.Logs = append(upload.Logs, models.Log{Message: "Generated file URL successfully", TimeStamp: now})
		upload.Logs = append(upload.Logs, models.Log{Message: "File upload completed successfully", TimeStamp: now})
	}

	return us.upRepo.Update(ctx, upload.WorkspaceID.Hex(), upload.ID.Hex(), upload)
}

func (us *UploadService) getUploaderConfig(hook *tusd.HookEvent) (*models.UploaderConfig, error) {
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return nil, err
	}

	config, err := us.wsRepo.GetUploaderConfig(hook.Context, workspaceID)
	if err != nil {
		return nil, fmt.Errorf(messages.WorkspaceNotFound, workspaceID)
	}

	return config, nil
}
