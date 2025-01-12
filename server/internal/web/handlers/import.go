package handlers

import (
	"net/http"

	"github.com/shivamsanju/uploader/internal/db/repo"
	"github.com/shivamsanju/uploader/pkg/globals"
	"github.com/shivamsanju/uploader/pkg/hooks"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

type importHandler struct {
	upRepo repo.UploaderRepo
}

func NewImportHandler() *importHandler {
	return &importHandler{
		upRepo: repo.NewUploaderRepo(),
	}
}

func (h *importHandler) GetTusHandler() http.Handler {
	globals.Log.Infof("initializing tusd handler with upload dir: %s", globals.TusUploadDir)
	store := filestore.New(globals.TusUploadDir)
	locker := filelocker.New(globals.TusUploadDir)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new tusd handler
	tusdHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:      "/imports",
		StoreComposer: composer,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			globals.Log.Infof("pre upload create -> %s", hook.HTTPRequest.URI)
			return tusd.HTTPResponse{
				StatusCode: http.StatusOK,
			}, tusd.FileInfoChanges{}, nil
		},
		PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
			globals.Log.Infof("pre finish response -> %s", hook.Upload.ID)
			defer hooks.RemoveTusUploadHook(hook)

			err := hooks.UploadToDatastoreHook(hook)
			if err != nil {
				globals.Log.Errorf("unable to upload to datastore: %s", err)
				return tusd.HTTPResponse{
					StatusCode: http.StatusBadRequest,
				}, nil
			}
			return tusd.HTTPResponse{
				StatusCode: http.StatusOK,
			}, nil
		},
	})

	if err != nil {
		globals.Log.Errorf("unable to create tusd handler: %s", err)
		panic(err)
	}
	return tusdHandler
}
