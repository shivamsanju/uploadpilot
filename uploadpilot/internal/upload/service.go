package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadService struct {
	upRepo      *db.UploadRepo
	wsRepo      *db.WorkspaceRepo
	logRepo     *db.UploadLogsRepo
	eventBus    *events.UploadEventBus
	logEventBus *events.LogEventBus
}

func NewUploadService() *UploadService {
	return &UploadService{
		upRepo:      db.NewUploadRepo(),
		wsRepo:      db.NewWorkspaceRepo(),
		logRepo:     db.NewUploadLogsRepo(),
		eventBus:    events.GetUploadEventBus(),
		logEventBus: events.GetLogEventBus(),
	}
}

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
	us.logEventBus.Publish(events.NewLogEvent(ctx, workspaceID, uploadID, "uploader details fetched", models.UploadLogLevelInfo))
	return us.upRepo.Get(ctx, uploadID)
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
		ID:        primitive.NewObjectID(),
		Size:      hook.Upload.Size,
		Status:    models.UploadStatusInProgress,
		StartedAt: now,
	}
	us.logEventBus.Publish(events.NewLogEvent(hook.Context, workspaceID, upload.ID.Hex(), "upload started", models.UploadLogLevelInfo))

	logErrorAndUpdateUpload := func(err error) error {
		if repoErr := us.upRepo.Create(hook.Context, workspaceID, upload); repoErr != nil {
			infra.Log.Errorf("unable to create upload: %s", repoErr)
			return repoErr
		}

		us.eventBus.Publish(events.NewUploadEvent(hook.Context, events.EventUploadFailed, upload, "", err))
		us.logEventBus.Publish(events.NewLogEvent(hook.Context, workspaceID, upload.ID.Hex(), "upload failed: "+err.Error(), models.UploadLogLevelError))

		return err
	}

	metadata, err := us.extractMetadataFromTusdEvent(hook)
	if err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}
	upload.Metadata = metadata

	if err := us.ValidateUploadSizeLimits(hook, workspaceID, upload.ID.Hex(), config); err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}

	if err := us.ValidateUploadFileType(hook, workspaceID, upload.ID.Hex(), config); err != nil {
		return nil, logErrorAndUpdateUpload(err)
	}

	if err := us.AuthenticateUpload(hook, workspaceID, upload.ID.Hex(), config); err != nil {
		return nil, logErrorAndUpdateUpload(fmt.Errorf(msg.UploadAuthenticationFailed, err))
	}

	if err := us.upRepo.Create(hook.Context, workspaceID, upload); err != nil {
		infra.Log.Errorf("unable to create upload: %s", err)
		return nil, err
	}
	us.eventBus.Publish(events.NewUploadEvent(hook.Context, events.EventUploadStarted, upload, "", nil))

	return upload, nil
}

func (us *UploadService) FinishUpload(ctx context.Context, uploadID string) error {
	upload, err := us.upRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	infra.Log.Infof("upload completed: %s", uploadID)
	us.eventBus.Publish(events.NewUploadEvent(ctx, events.EventUploadComplete, upload, "", nil))
	us.logEventBus.Publish(events.NewLogEvent(ctx, upload.WorkspaceID.Hex(), upload.ID.Hex(), "upload completed", models.UploadLogLevelInfo))
	return nil
}

func (us *UploadService) GetLogs(ctx context.Context, uploadID string) ([]bson.M, error) {
	logs, err := us.logRepo.GetLogs(ctx, uploadID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (us *UploadService) getUploaderConfig(hook *tusd.HookEvent) (*models.UploaderConfig, error) {
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return nil, err
	}

	config, err := us.wsRepo.GetUploaderConfig(hook.Context, workspaceID)
	if err != nil {
		return nil, fmt.Errorf(msg.WorkspaceNotFound, workspaceID)
	}

	return config, nil
}

func (us *UploadService) extractMetadataFromTusdEvent(hook *tusd.HookEvent) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
	if err != nil {
		infra.Log.Errorf("Failed to extract metadata from upload request: %s", err.Error())
		return metadata, err
	}
	return metadata, nil
}

func (us *UploadService) getWorkspaceIDFromTusdEvent(hook *tusd.HookEvent) (string, error) {
	infra.Log.Infof("Upload Object %+v", hook.Upload)

	headers := hook.HTTPRequest.Header
	workspaceID := headers.Get("workspaceId")
	if len(workspaceID) == 0 {
		return "", fmt.Errorf(msg.InvalidWorkspaceIDInHeaders, workspaceID)
	}

	return workspaceID, nil
}
