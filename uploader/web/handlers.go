package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploader/internal/service"
	"golang.org/x/exp/slog"
)

type handler struct {
	uploadSvc *service.Service
}

func Newhandler(uploadSvc *service.Service) *handler {
	return &handler{
		uploadSvc: uploadSvc,
	}
}

func (h *handler) HealthCheckWithCompanion(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(config.GetAppConfig().CompanionEndpoint + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		HandleHttpError(w, r, http.StatusServiceUnavailable, errors.New("companion is not healthy"))
		return
	}
	render.JSON(w, r, "uploader is healthy")
}

func (h *handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "uploader is healthy")
}

func (h *handler) GetUploaderConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	err := VerifyWorkspaceID(workspaceID)
	if err != nil {
		HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	cfg, err := h.uploadSvc.GetUploaderConfig(r.Context(), r, workspaceID)
	if err != nil {
		HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, cfg)
}

func (h *handler) GetUploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workspaceID := chi.URLParam(r, "workspaceId")
		err := VerifyWorkspaceID(workspaceID)
		if err != nil {
			HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}

		handlerConfig, err := h.uploadSvc.GetTusdConfigForWorkspace(r.Context(), r, workspaceID)
		if err != nil {
			HandleHttpError(w, r, http.StatusBadRequest, err)
			return
		}
		handlerConfig.Logger = slog.New(config.NewLogHandler())
		// tusd handler
		uploadHandler, err := tusd.NewHandler(*handlerConfig)
		if err != nil {
			// infra.Log.Errorf("unable to create tusd handler: %s", err)
			HandleHttpError(w, r, http.StatusInternalServerError, errors.New("unknown error"))
			return
		}

		http.StripPrefix(fmt.Sprintf("/upload/%s", workspaceID), uploadHandler).ServeHTTP(w, r)
	})

}
