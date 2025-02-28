package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/workspace"
	"github.com/uploadpilot/manager/internal/utils"
)

type workspaceHandler struct {
	workspaceSvc *workspace.Service
}

func NewWorkspaceHandler(workspaceSvc *workspace.Service) *workspaceHandler {
	return &workspaceHandler{
		workspaceSvc: workspaceSvc,
	}
}

func (h *workspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace := &models.Workspace{}
	if err := render.DecodeJSON(r.Body, workspace); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.CreateWorkspace(r.Context(), workspace)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, workspace.ID)
}

func (h *workspaceHandler) GetAllWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.workspaceSvc.GetAllWorkspaces(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, workspaces)
}

func (h *workspaceHandler) GetWorkspaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	uploaderConfig, err := h.workspaceSvc.GetWorkspaceConfig(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, uploaderConfig)
}

func (h *workspaceHandler) SetWorkspaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	uploaderConfig := &models.WorkspaceConfig{}
	if err := render.DecodeJSON(r.Body, uploaderConfig); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.SetWorkspaceConfig(r.Context(), workspaceID, uploaderConfig)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
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
