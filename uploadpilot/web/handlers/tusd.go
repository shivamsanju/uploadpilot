package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
	"github.com/uploadpilot/uploadpilot/internal/config"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/hooks"
	"github.com/uploadpilot/uploadpilot/internal/hooks/catalog"
	"github.com/uploadpilot/uploadpilot/internal/storage"

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
	s3Client, err := storage.S3Client(context.Background())
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
			// Accept the upload
			infra.Log.Infof("creating import for %s", hook.Upload.ID)
			err := h.validateWorkspaceId(&hook)
			if err != nil {
				return tusdBadRequestResponse(), tusd.FileInfoChanges{}, err
			}

			return tusdOkResponse(), tusd.FileInfoChanges{}, nil
		},

		PreFinishResponseCallback: func(tusdHook tusd.HookEvent) (tusd.HTTPResponse, error) {
			imp, err := h.createImport(&tusdHook)
			if err != nil {
				infra.Log.Errorf("unable to get import: %s", err)
				return tusdBadRequestResponse(), nil
			}
			prefinishExecutor := catalog.BuildPrefinishResponseHookExecutor()
			err = prefinishExecutor.Start(tusdHook.Context, &hooks.HookInput{
				Import:   imp,
				TusdHook: &tusdHook,
			}, false)
			if err != nil {
				infra.Log.Errorf("unable to execute prefinish hook: %s", err)
				return tusdBadRequestResponse(), nil
			}

			postfinishExecutor := catalog.BuildPostfinishResponseHookExecutor()
			go postfinishExecutor.Start(tusdHook.Context, &hooks.HookInput{
				Import:   imp,
				TusdHook: &tusdHook,
			}, false)

			infra.Log.Infof("tusd upload finished for %s", tusdHook.Upload.ID)

			return tusdOkResponse(), nil
		},
	})

	if err != nil {
		infra.Log.Errorf("unable to create tusd handler: %s", err)
		panic(err)
	}
	return tusdHandler
}

func (tusdHandler *tusdHandler) validateWorkspaceId(tusdHook *tusd.HookEvent) error {
	headers := tusdHook.HTTPRequest.Header
	wsID := headers.Get("workspaceId")
	if len(wsID) == 0 {
		return errors.New("missing workspaceId in header")
	}

	workspaceID, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		return fmt.Errorf("invalid workspaceId: %w", err)
	}

	wsRepo := db.NewWorkspaceRepo()
	exists, err := wsRepo.CheckWorkspaceExists(tusdHook.Context, workspaceID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("workspace not found")
	}

	return nil
}

func (h *tusdHandler) createImport(tusdHook *tusd.HookEvent) (*models.Import, error) {
	infra.Log.Infof("creating import for %s", tusdHook.Upload.ID)
	headers := tusdHook.HTTPRequest.Header
	wsID := headers.Get("workspaceId")
	if len(wsID) == 0 {
		return nil, errors.New("missing workspaceId in header")
	}

	workspaceID, err := primitive.ObjectIDFromHex(wsID)
	if err != nil {
		return nil, fmt.Errorf("invalid workspaceId: %w", err)
	}

	imp := models.Import{}

	imp.ID = primitive.NewObjectID()
	imp.UploadID = tusdHook.Upload.ID
	imp.Size = tusdHook.Upload.Size
	imp.WorkspaceID = workspaceID
	imp.StartedAt = primitive.NewDateTimeFromTime(time.Now())
	imp.Status = models.ImportStatusInProgress
	imp.Logs = []models.Log{{
		Message:   "Import created successfully",
		TimeStamp: primitive.NewDateTimeFromTime(time.Now()),
	}}

	createdImport, err := h.impRepo.Create(tusdHook.Context, &imp)
	return createdImport, err
}

func tusdOkResponse() tusd.HTTPResponse {
	return tusd.HTTPResponse{
		StatusCode: http.StatusOK,
	}
}

func tusdBadRequestResponse(statusCode ...int) tusd.HTTPResponse {
	if len(statusCode) > 0 && statusCode[0] > 0 {
		return tusd.HTTPResponse{
			StatusCode: statusCode[0],
		}
	}
	return tusd.HTTPResponse{
		StatusCode: http.StatusBadRequest,
	}
}
