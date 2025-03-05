package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/web/webutils"
)

type workspaceHandler struct {
	workspaceSvc *services.WorkspaceService
}

func NewWorkspaceHandler(workspaceSvc *services.WorkspaceService) *workspaceHandler {
	return &workspaceHandler{
		workspaceSvc: workspaceSvc,
	}
}

func (h *workspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace := &models.Workspace{}
	if err := render.DecodeJSON(r.Body, workspace); err != nil {
		webutils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.CreateWorkspace(r.Context(), workspace)
	if err != nil {
		webutils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, workspace.ID)
}

func (h *workspaceHandler) GetAllWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.workspaceSvc.GetAllWorkspaces(r.Context())
	if err != nil {
		webutils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, workspaces)
}

func (h *workspaceHandler) GetWorkspaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	uploaderConfig, err := h.workspaceSvc.GetWorkspaceConfig(r.Context(), workspaceID)
	if err != nil {
		webutils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, uploaderConfig)
}

func (h *workspaceHandler) SetWorkspaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	uploaderConfig := &models.WorkspaceConfig{}
	if err := render.DecodeJSON(r.Body, uploaderConfig); err != nil {
		webutils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.SetWorkspaceConfig(r.Context(), workspaceID, uploaderConfig)
	if err != nil {
		webutils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (h *workspaceHandler) GetAllAllowedSources(
	r *http.Request, params dto.WorkspaceParams, query interface{}, body interface{},
) ([]string, int, error) {
	return models.AllAllowedSources, http.StatusOK, nil
}

func (h *workspaceHandler) LogUpload(
	r *http.Request, params dto.WorkspaceParams, query interface{}, body interface{},
) (string, int, error) {
	return "", http.StatusOK, nil
}
