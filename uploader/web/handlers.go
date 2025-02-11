package web

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/dto"
	uploaderconfig "github.com/uploadpilot/uploadpilot/uploader/internal/svc/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/upload"

	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
)

type handler struct {
	uploadSvc upload.UploadService
	configSvc uploaderconfig.UploaderConfigService
}

func Newhandler() *handler {
	return &handler{
		uploadSvc: *upload.NewUploadService(),
		configSvc: *uploaderconfig.NewUploaderConfigService(),
	}
}

func (h *handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(config.CompanionEndpoint + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		HandleHttpError(w, r, http.StatusServiceUnavailable, errors.New("companion is not healthy"))
		return
	}

	render.JSON(w, r, "uploader is healthy")
}

func (h *handler) TusHandler() http.Handler {
	infra.Log.Infof("initializing tusd handler with upload dir: %s", config.TusUploadDir)

	// S3 backend for tusd
	s3Client, err := infra.NewS3Client(&infra.S3Config{
		AccessKey: config.S3AccessKey,
		SecretKey: config.S3SecretKey,
		Region:    config.S3Region,
	})
	if err != nil {
		panic(err)
	}
	store := s3store.New(config.S3BucketName, s3Client)
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	// tusd handler
	infra.Log.Infof("initializing tusd handler with upload base path: %s", config.TusUploadBasePath)
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:           config.TusUploadBasePath,
		StoreComposer:      composer,
		MaxSize:            config.TusMaxFileSize,
		DisableDownload:    true,
		DisableTermination: false,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			active, err := h.uploadSvc.VerifySubscription(&hook)
			if err != nil {
				infra.Log.Errorf("unable to verify subscription: %s", err)
				return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, nil
			}
			if !active {
				return tusd.HTTPResponse{StatusCode: http.StatusForbidden, Body: "subscription is not active"}, tusd.FileInfoChanges{}, nil
			}

			upload, err := h.uploadSvc.CreateUpload(&hook)
			if err != nil {
				infra.Log.Errorf("unable to create upload: %s", err)
				return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, nil
			}

			hook.Upload.MetaData["upload_id"] = upload.ID
			return tusd.HTTPResponse{StatusCode: http.StatusOK}, tusd.FileInfoChanges{
				ID:       upload.ID,
				MetaData: hook.Upload.MetaData,
			}, nil
		},

		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			uploadID := hook.Upload.MetaData["upload_id"]
			infra.Log.Infof("trying to finish upload: %s", uploadID)
			if err := h.uploadSvc.FinishUpload(hook.Context, uploadID); err != nil {
				infra.Log.Errorf("unable to finish upload: %s", err)
				return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, nil
			}

			return tusd.HTTPResponse{StatusCode: http.StatusOK}, nil
		},
	})

	if err != nil {
		infra.Log.Errorf("unable to create tusd handler: %s", err)
		panic(err)
	}
	return tusdHandler
}

func (h *handler) GetUploaderConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	uploaderConfig, err := h.configSvc.GetUploaderConfig(r.Context(), workspaceID)
	if err != nil {
		HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	cfg := dto.UploaderConfig{
		UploaderConfig: *uploaderConfig,
		ChunkSize:      config.TusChunkSize,
	}
	render.JSON(w, r, cfg)
}
