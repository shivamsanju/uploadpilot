package handlers

import (
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"

	"github.com/uploadpilot/uploadpilot/internal/hooks"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tusdHandler struct {
	impRepo db.ImportRepo
	wsRepo  db.WorkspaceRepo
}

func NewTusdHandler() *tusdHandler {
	return &tusdHandler{
		impRepo: db.NewImportRepo(),
		wsRepo:  db.NewWorkspaceRepo(),
	}
}

func (h *tusdHandler) GetTusHandler() http.Handler {
	infra.Log.Infof("initializing tusd handler with upload dir: %s", config.TusUploadDir)
	store := filestore.New(config.TusUploadDir)
	locker := filelocker.New(config.TusUploadDir)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new tusd handler
	infra.Log.Infof("initializing tusd handler with upload base path: %s", config.TusUploadBasePath)
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:                config.TusUploadBasePath,
		StoreComposer:           composer,
		RespectForwardedHeaders: true,
		// NotifyTerminatedUploads: true,
		// NotifyCompleteUploads:   true,
		// NotifyCreatedUploads:    true,
		// NotifyUploadProgress:    true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			infra.Log.Infof("pre upload create -> %s", hook.HTTPRequest.URI)

			// Validate the upload
			err := hooks.ValidateUpload(hook)
			if err != nil {
				infra.Log.Errorf("unable to validate upload: %s", err)
				return tusdBadRequestResponse(), tusd.FileInfoChanges{}, nil
			}

			return tusdOkResponse(), tusd.FileInfoChanges{}, nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			infra.Log.Infof("pre finish response -> %s", hook.Upload.ID)

			// Extract Metadata
			var metadata map[string]interface{}
			err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
			if err != nil {
				return tusdBadRequestResponse(), nil
			}

			// Remove uploads from local
			fileName, ok := metadata["filename"]
			if !ok || len(fileName.(string)) == 0 {
				return tusdBadRequestResponse(), nil
			}
			defer hooks.RenameTusUploadHook(hook, fileName.(string))

			// Create new upload
			importId := primitive.NewObjectID()
			imp := &models.Import{
				ID:        importId,
				Size:      hook.Upload.Size,
				StartedAt: primitive.NewDateTimeFromTime(time.Now()),
				Metadata:  metadata,
				Logs: []models.Log{{
					Message:   "Import started",
					TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
				}},
				Status: models.ImportStatusInProgress,
			}
			go h.impRepo.Create(hook.Context, imp)

			// Update metadata
			err = hooks.UpdateImportMetadata(hook, imp)
			if err != nil {
				imp.Logs = append(imp.Logs, models.Log{
					Message:   "Error updating the metadata of the file: " + err.Error(),
					TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
				})
				imp.Status = models.ImportStatusFailed
				infra.Log.Errorf("unable to upload to datastore: %s", err)
				h.impRepo.Update(hook.Context, importId, imp)
				return tusdOkResponse(), nil
			}

			// Update upload
			imp.Logs = append(imp.Logs, models.Log{
				Message:   "Import completed",
				TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
			})
			imp.Status = models.ImportStatusSuccess
			_, err = h.impRepo.Update(hook.Context, importId, imp)
			if err != nil {
				infra.Log.Errorf("unable to update import: %s", err)
				return tusdOkResponse(), nil
			}

			return tusdOkResponse(), nil
		},
	})

	if err != nil {
		infra.Log.Errorf("unable to create tusd handler: %s", err)
		panic(err)
	}
	return tusdHandler
}

func tusdOkResponse() tusd.HTTPResponse {
	return tusd.HTTPResponse{
		StatusCode: http.StatusOK,
	}
}

func tusdBadRequestResponse() tusd.HTTPResponse {
	return tusd.HTTPResponse{
		StatusCode: http.StatusBadRequest,
	}
}
