package handlers

import (
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	"github.com/uploadpilot/uploadpilot/pkg/db/repo"
	"github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/hooks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tusdHandler struct {
	upRepo  repo.UploaderRepo
	impRepo repo.ImportRepo
}

func NewTusdHandler() *tusdHandler {
	return &tusdHandler{
		upRepo:  repo.NewUploaderRepo(),
		impRepo: repo.NewImportRepo(),
	}
}

func (h *tusdHandler) GetTusHandler() http.Handler {
	globals.Log.Infof("initializing tusd handler with upload dir: %s", globals.TusUploadDir)
	store := filestore.New(globals.TusUploadDir)
	locker := filelocker.New(globals.TusUploadDir)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new tusd handler
	globals.Log.Infof("initializing tusd handler with upload base path: %s", globals.TusUploadBasePath)
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:      globals.TusUploadBasePath,
		StoreComposer: composer,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			globals.Log.Infof("pre upload create -> %s", hook.HTTPRequest.URI)

			// Validate the upload
			err := hooks.ValidateUpload(hook)
			if err != nil {
				globals.Log.Errorf("unable to validate upload: %s", err)
				return tusdBadRequestResponse(), tusd.FileInfoChanges{}, nil
			}

			return tusdOkResponse(), tusd.FileInfoChanges{}, nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			globals.Log.Infof("pre finish response -> %s", hook.Upload.ID)

			// Remove uploads from local
			defer hooks.RemoveTusUploadHook(hook)

			// Extract Metadata
			var metadata map[string]interface{}
			err := mapstructure.Decode(hook.Upload.MetaData, &metadata)
			if err != nil {
				return tusdBadRequestResponse(), nil
			}

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

			// Upload to datastore
			err = hooks.UploadToDatastoreHook(hook, imp)
			if err != nil {
				imp.Logs = append(imp.Logs, models.Log{
					Message:   "Error importing the file: " + err.Error(),
					TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
				})
				imp.Status = models.ImportStatusFailed
				globals.Log.Errorf("unable to upload to datastore: %s", err)
				h.impRepo.Update(hook.Context, &importId, imp)
				return tusdOkResponse(), nil
			}

			// Update upload
			imp.Logs = append(imp.Logs, models.Log{
				Message:   "Import completed",
				TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
			})
			imp.Status = models.ImportStatusSuccess
			_, err = h.impRepo.Update(hook.Context, &importId, imp)
			if err != nil {
				globals.Log.Errorf("unable to update import: %s", err)
				return tusdOkResponse(), nil
			}

			return tusdOkResponse(), nil
		},
	})

	if err != nil {
		globals.Log.Errorf("unable to create tusd handler: %s", err)
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
