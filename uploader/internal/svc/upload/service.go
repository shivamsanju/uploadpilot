package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/common/pkg/db/repo"
	dbutils "github.com/uploadpilot/uploadpilot/common/pkg/db/utils"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/msg"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	uploaderconfig "github.com/uploadpilot/uploadpilot/uploader/internal/svc/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/workspace"
	"github.com/uploadpilot/uploadpilot/uploader/internal/validations"
)

type Service struct {
	uploadRepo     *repo.UploadRepo
	uploadLogsRepo *repo.UploadLogsRepo
	wsSvc          *workspace.Service
	configSvc      *uploaderconfig.Service
	logEventBus    *pubsub.EventBus[events.UploadLogEventMsg]
	uploadEb       *pubsub.EventBus[events.UploadEventMsg]
}

func NewUploadService(uploadRepo *repo.UploadRepo, uploadLogsRepo *repo.UploadLogsRepo,
	workspaceRepo *repo.WorkspaceRepo, configRepo *repo.WorkspaceConfigRepo) *Service {
	return &Service{
		uploadRepo:     uploadRepo,
		uploadLogsRepo: uploadLogsRepo,
		wsSvc:          workspace.NewWorkspaceService(workspaceRepo),
		configSvc:      uploaderconfig.NewConfigService(configRepo),
		uploadEb:       events.NewUploadStatusEvent(infra.RedisClient, uuid.New().String()),
		logEventBus:    events.NewUploadLogEventBus(infra.RedisClient, uuid.New().String()),
	}
}

func (us *Service) VerifySubscription(hook *tusd.HookEvent) (bool, error) {
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return false, errors.New("invalid workspace id in headers")
	}
	active, err := us.wsSvc.VerifySubscription(hook.Context, workspaceID)
	return active, err
}

func (us *Service) CreateUpload(hook *tusd.HookEvent) (*models.Upload, error) {
	workspaceID, err := us.getWorkspaceIDFromTusdEvent(hook)
	if err != nil {
		return nil, errors.New("invalid workspace id in headers")
	}

	// workspace existence is anyways checked here, no need to add extra method call
	config, err := us.configSvc.GetUploaderConfig(hook.Context, workspaceID)
	if err != nil {
		return nil, err
	}

	id := dbutils.GenerateUUID()
	upload := &models.Upload{
		ID:          id,
		Size:        hook.Upload.Size,
		Status:      models.UploadStatusInProgress,
		WorkspaceID: workspaceID,
		Metadata:    map[string]interface{}{},
	}

	eventMsg := events.NewUploadEventMessage(workspaceID, upload.ID, string(models.UploadStatusInProgress), nil, nil)

	metadata, err := us.extractMetadataFromTusdEvent(hook)
	if err != nil {
		upload.Status = models.UploadStatusFailed
		if err := us.uploadRepo.Create(hook.Context, workspaceID, upload); err != nil {
			infra.Log.Errorf(msg.FailedToCreateUpload, err.Error())
			return nil, err
		}
		eventMsg.Status = string(models.UploadStatusFailed)
		us.uploadEb.Publish(eventMsg)
		return nil, err
	}

	metadata["upload_id"] = id
	upload.Metadata = metadata
	if err := us.uploadRepo.Create(hook.Context, workspaceID, upload); err != nil {
		infra.Log.Errorf("unable to create upload: %s", err)
		return nil, err
	}
	us.uploadEb.Publish(eventMsg)

	validators := []func(*tusd.HookEvent, string, string, *models.UploaderConfig) (string, error){
		validations.ValidateUploadSizeLimits,
		validations.ValidateUploadFileType,
		validations.AuthenticateUpload,
	}

	for _, validator := range validators {
		if err := us.processValidation(hook, eventMsg, validator, config); err != nil {
			return nil, err
		}
	}

	return upload, nil
}

func (us *Service) FinishUpload(ctx context.Context, uploadID string) error {
	upload, err := us.uploadRepo.Get(ctx, uploadID)
	if err != nil {
		return err
	}
	infra.Log.Infof(msg.UploadComplete + ": " + uploadID)
	eventMsg := events.NewUploadEventMessage(upload.WorkspaceID, upload.ID, string(models.UploadStatusComplete), nil, nil)
	us.uploadEb.Publish(eventMsg)
	return nil
}

func (us *Service) SetStatus(ctx context.Context, uploadID string, status models.UploadStatus) error {
	return us.uploadRepo.SetStatus(ctx, uploadID, status)
}

func (us *Service) BatchAddLogs(ctx context.Context, logs []*models.UploadLog) error {
	return us.uploadLogsRepo.BatchAddLogs(ctx, logs)
}

func (us *Service) getWorkspaceIDFromTusdEvent(hook *tusd.HookEvent) (string, error) {
	infra.Log.Infof("Upload Object %+v", hook.Upload)

	headers := hook.HTTPRequest.Header
	workspaceID := headers.Get("workspaceId")
	if len(workspaceID) == 0 {
		return "", fmt.Errorf(msg.InvalidWorkspaceIDInHeaders, workspaceID)
	}

	return workspaceID, nil
}

func (us *Service) extractMetadataFromTusdEvent(hook *tusd.HookEvent) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
	if err != nil {
		infra.Log.Errorf("Failed to extract metadata from upload request: %s", err.Error())
		return metadata, err
	}
	return metadata, nil
}

func (us *Service) processValidation(
	hook *tusd.HookEvent,
	eventMsg *events.UploadEventMsg,
	validator func(*tusd.HookEvent, string, string, *models.UploaderConfig) (string, error),
	config *models.UploaderConfig,
) error {
	success, err := validator(hook, eventMsg.WorkspaceID, eventMsg.UploadID, config)
	if err != nil {
		eventMsg.Status = string(models.UploadStatusFailed)
		us.uploadEb.Publish(eventMsg)
		return err
	}
	us.logEventBus.Publish(events.NewUploadLogEventMessage(
		eventMsg.WorkspaceID, eventMsg.UploadID, nil, nil, success, models.UploadLogLevelInfo,
	))
	return nil
}
