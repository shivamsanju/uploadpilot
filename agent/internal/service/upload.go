package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/phuslu/log"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	"github.com/uploadpilot/agent/internal/clients"
	"github.com/uploadpilot/agent/internal/config"
	"github.com/uploadpilot/agent/internal/dto"
)

type Service struct {
	coreSvcClient *clients.CoreServiceClient
}

func NewUploadService(c *clients.CoreServiceClient) *Service {
	return &Service{
		coreSvcClient: c,
	}
}

func (us *Service) GetUploaderConfig(ctx context.Context, originalReq *http.Request, workspaceID string) (*dto.WorkspaceConfig, int, error) {
	cgf, statusCode, err := us.coreSvcClient.GetUploaderConfig(ctx, originalReq, workspaceID)
	if err != nil {
		log.Error().Err(err).Str("workspace_id", workspaceID).Msg("failed to get uploader config")
		return nil, statusCode, err
	}
	return cgf, statusCode, nil
}

func (us *Service) GetTusdConfigForWorkspace(ctx context.Context, originalReq *http.Request, workspaceID string) (*tusd.Config, int, error) {
	appConfig := config.GetAppConfig()

	uploaderConfig, statusCode, err := us.GetUploaderConfig(ctx, originalReq, workspaceID)
	if err != nil {
		log.Error().Err(err).Str("workspace_id", workspaceID).Msg("failed to get uploader config")
		return nil, statusCode, err
	}

	corsCfg, err := us.getCorsConfigForWorkspace(uploaderConfig)
	if err != nil {
		log.Error().Err(err).Str("workspace_id", workspaceID).Msg("failed to get cors config")
		return nil, statusCode, err
	}

	composer, err := us.getStorageForWorkspace()
	if err != nil {
		log.Error().Err(err).Str("workspace_id", workspaceID).Msg("failed to get storage for workspace")
		return nil, statusCode, err
	}

	handlerConfig := tusd.Config{
		BasePath:                  fmt.Sprintf(appConfig.TusUploadBasePath, workspaceID),
		StoreComposer:             composer,
		DisableDownload:           true,
		DisableTermination:        false,
		PreUploadCreateCallback:   us.getPreUploadCallback(originalReq, workspaceID, uploaderConfig),
		PreFinishResponseCallback: us.getPreFinishCallback(originalReq, workspaceID),
		Cors:                      corsCfg,
	}

	// max file size
	if uploaderConfig.MaxFileSize != nil {
		handlerConfig.MaxSize = *uploaderConfig.MaxFileSize
	}

	return &handlerConfig, statusCode, nil
}

// TUSD configs

func (us *Service) getStorageForWorkspace() (*tusd.StoreComposer, error) {
	cfg := config.GetAppConfig()
	s3Client, err := clients.NewS3Client(&clients.S3Options{
		AccessKey: cfg.S3AccessKey,
		SecretKey: cfg.S3SecretKey,
		Region:    cfg.S3Region,
	})
	if err != nil {
		return nil, err
	}

	store := s3store.New(cfg.S3BucketName, s3Client)
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	return composer, nil
}

func (us *Service) getCorsConfigForWorkspace(uploaderConfig *dto.WorkspaceConfig) (*tusd.CorsConfig, error) {
	cf := &tusd.CorsConfig{
		AllowCredentials: false,
		Disable:          false,
		AllowMethods:     "PATCH, POST, HEAD, OPTIONS",
		AllowHeaders:     tusd.DefaultCorsConfig.AllowHeaders,
		MaxAge:           tusd.DefaultCorsConfig.MaxAge,
		ExposeHeaders:    tusd.DefaultCorsConfig.ExposeHeaders,
	}

	if len(uploaderConfig.AllowedOrigins) > 0 {
		// Escape and join origins properly for regex matching
		escapedOrigins := make([]string, len(uploaderConfig.AllowedOrigins))
		for i, origin := range uploaderConfig.AllowedOrigins {
			escapedOrigins[i] = regexp.QuoteMeta(origin) // Escape special regex chars
		}

		pattern := "^(" + strings.Join(escapedOrigins, "|") + ")$"
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return cf, err
		}
		cf.AllowOrigin = regex
	} else {
		cf.AllowOrigin = tusd.DefaultCorsConfig.AllowOrigin
	}

	return cf, nil

}

func (us *Service) extractMetadataFromTusdEvent(hook *tusd.HookEvent) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

func (us *Service) logUploadRequest(hook *tusd.HookEvent, originalReq *http.Request, workspaceID string) (int, error) {
	data := map[string]string{}
	for k, v := range hook.HTTPRequest.Header {
		data[k] = strings.Join(v, ",")
	}

	return us.coreSvcClient.LogUploadRequest(hook.Context, originalReq, workspaceID, data)
}

func (us *Service) getFileNameFromMetadata(metadata map[string]interface{}) (string, error) {
	filename, ok := metadata["filename"].(string)
	if !ok {
		return "", fmt.Errorf("filename not found in metadata")
	}
	return filename, nil
}

func (us *Service) getFileTypeFromMetadata(metadata map[string]interface{}) (string, error) {
	fileType, ok := metadata["filetype"].(string)
	if !ok {
		return "", fmt.Errorf("filetype not found in metadata")
	}
	return fileType, nil
}

func (us *Service) getPreUploadCallback(originalReq *http.Request, workspaceID string, uploaderConfig *dto.WorkspaceConfig) func(tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
	return func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
		if statusCode, err := us.logUploadRequest(&hook, originalReq, workspaceID); err != nil {
			log.Error().Str("workspace_id", workspaceID).Err(err).Msg("unable to log upload request")
			return tusd.HTTPResponse{StatusCode: statusCode}, tusd.FileInfoChanges{}, errors.New("some error occurred, please try again")
		}
		metadata, err := us.extractMetadataFromTusdEvent(&hook)
		if err != nil {
			log.Error().Str("workspace_id", workspaceID).Err(err).Msg("unable to extract metadata from tusd event")
			return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, errors.New("unable to extract metadata")
		}

		fileName, err := us.getFileNameFromMetadata(metadata)
		if err != nil {
			log.Error().Str("workspace_id", workspaceID).Err(err).Msg("unable to extract filename from metadata")
			return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, errors.New("unable to extract filename")
		}

		fileType, err := us.getFileTypeFromMetadata(metadata)
		if err != nil {
			log.Error().Str("workspace_id", workspaceID).Err(err).Msg("unable to extract filetype from metadata")
			return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, errors.New("unable to extract filetype")
		}

		if len(uploaderConfig.AllowedFileTypes) > 0 && !slices.Contains(uploaderConfig.AllowedFileTypes, fileType) {
			log.Error().Str("workspace_id", workspaceID).Str("filetype", fileType).Msg("filetype not allowed")
			return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, errors.New("filetype not allowed")
		}

		upload := &dto.Upload{
			ID:        uuid.New().String(),
			Size:      hook.Upload.Size,
			Metadata:  metadata,
			FileName:  fileName,
			FileType:  fileType,
			StartedAt: time.Now(),
		}

		uploadID, statusCode, err := us.coreSvcClient.CreateNewUpload(hook.Context, originalReq, workspaceID, upload)
		if err != nil {
			log.Error().Str("workspace_id", workspaceID).Err(err).Msg("unable to create new upload")
			return tusd.HTTPResponse{StatusCode: statusCode}, tusd.FileInfoChanges{}, errors.New("some error occurred, please try again")
		}

		hook.Upload.MetaData["upload_id"] = upload.ID
		return tusd.HTTPResponse{StatusCode: http.StatusCreated}, tusd.FileInfoChanges{
			ID:       uploadID,
			MetaData: hook.Upload.MetaData,
		}, nil
	}
}

func (us *Service) getPreFinishCallback(originalReq *http.Request, workspaceID string) func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
	return func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
		uploadID := hook.Upload.MetaData["upload_id"]
		if uploadID == "" {
			hook.Upload.StopUpload(tusd.HTTPResponse{StatusCode: http.StatusBadRequest})
			log.Error().Str("workspace_id", workspaceID).Msg("failed to finish upload: upload id not found in metadata")
			return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, nil
		}

		if statusCode, err := us.coreSvcClient.FinishUpload(hook.Context, originalReq, workspaceID, uploadID, &dto.Upload{
			Status:     "Uploaded",
			FinishedAt: time.Now(),
		}); err != nil {
			hook.Upload.StopUpload(tusd.HTTPResponse{StatusCode: statusCode})
			log.Error().Str("workspace_id", workspaceID).Str("upload_id", uploadID).Err(err).Msg("unable to finish upload")
			return tusd.HTTPResponse{StatusCode: statusCode}, nil
		}
		return tusd.HTTPResponse{StatusCode: http.StatusNoContent}, nil
	}
}

func (us *Service) GetFileTypeFromTusdEvent(ctx context.Context, composer *tusd.StoreComposer, uploadID string) (string, error) {
	upload, err := composer.Core.GetUpload(ctx, uploadID)
	if err != nil {
		return "", err
	}
	reader, err := upload.GetReader(ctx)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		log.Error().Str("uploadID", uploadID).Err(err).Msg("unable to read file content")
		return "", err
	}

	// Detect content type
	contentType := http.DetectContentType(buf[:n])

	return contentType, nil
}
