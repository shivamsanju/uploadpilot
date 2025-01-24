package handlers

import (
	"net/http"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/upload"
	"github.com/uploadpilot/uploadpilot/internal/webhook"
	"github.com/uploadpilot/uploadpilot/internal/workspace"

	"github.com/uploadpilot/uploadpilot/internal/infra"
)

type tusdHandler struct {
	workspaceSvc *workspace.WorkspaceService
	uploadSvc    *upload.UploadService
	webhookSvc   *webhook.WebhookService
}

func NewTusdHandler() *tusdHandler {
	return &tusdHandler{
		workspaceSvc: workspace.NewWorkspaceService(),
		uploadSvc:    upload.NewUploadService(),
		webhookSvc:   webhook.NewWebhookService(),
	}
}

func (h *tusdHandler) GetTusHandler() http.Handler {
	infra.Log.Infof("initializing tusd handler with upload dir: %s", config.TusUploadDir)
	// A seperate s3 client for tusd
	s3Client, err := infra.NewS3Client()
	if err != nil {
		panic(err)
	}
	store := s3store.New(config.S3BucketName, s3Client)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	// Create a new tusd handler
	infra.Log.Infof("initializing tusd handler with upload base path: %s", config.TusUploadBasePath)
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:                config.TusUploadBasePath,
		StoreComposer:           composer,
		RespectForwardedHeaders: true,
		MaxSize:                 500 * 1024 * 1024,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			upload, err := h.uploadSvc.CreateUpload(&hook)
			if err != nil {
				infra.Log.Errorf("unable to create upload: %s", err)
				return tusd.HTTPResponse{StatusCode: http.StatusBadRequest}, tusd.FileInfoChanges{}, nil
			}

			hook.Upload.MetaData["uploadId"] = upload.ID.Hex()
			return tusd.HTTPResponse{StatusCode: http.StatusOK}, tusd.FileInfoChanges{
				ID:       upload.ID.Hex(),
				MetaData: hook.Upload.MetaData,
			}, nil
		},

		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			uploadID := hook.Upload.MetaData["uploadId"]
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
