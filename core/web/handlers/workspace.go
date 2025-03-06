package handlers

import (
	"net/http"

	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
)

type workspaceHandler struct {
	workspaceSvc *services.WorkspaceService
}

func NewWorkspaceHandler(workspaceSvc *services.WorkspaceService) *workspaceHandler {
	return &workspaceHandler{
		workspaceSvc: workspaceSvc,
	}
}

func (h *workspaceHandler) CreateWorkspace(r *http.Request, params dto.TenantParams, query interface{}, body models.Workspace) (string, int, error) {
	err := h.workspaceSvc.CreateWorkspace(r.Context(), params.TenantID, &body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return body.ID, http.StatusOK, nil
}

func (h *workspaceHandler) GetAllWorkspaces(r *http.Request, params dto.TenantParams, query, body interface{}) ([]models.Workspace, int, error) {
	workspaces, err := h.workspaceSvc.GetAllWorkspaces(r.Context(), params.TenantID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return workspaces, http.StatusOK, nil
}

func (h *workspaceHandler) GetWorkspaceConfig(r *http.Request, params dto.WorkspaceParams, query, body interface{}) (*models.WorkspaceConfig, int, error) {
	uploaderConfig, err := h.workspaceSvc.GetWorkspaceConfig(r.Context(), params.TenantID, params.WorkspaceID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return uploaderConfig, http.StatusOK, nil
}

func (h *workspaceHandler) SetWorkspaceConfig(r *http.Request, params dto.WorkspaceParams, query interface{}, body models.WorkspaceConfig) (bool, int, error) {
	err := h.workspaceSvc.SetWorkspaceConfig(r.Context(), params.TenantID, params.WorkspaceID, &body)
	if err != nil {
		return false, http.StatusInternalServerError, err
	}
	return true, http.StatusOK, nil
}

func (h *workspaceHandler) GetAllAllowedSources(r *http.Request, params dto.WorkspaceParams, query, body interface{}) ([]string, int, error) {
	return models.AllAllowedSources, http.StatusOK, nil
}

func (h *workspaceHandler) LogUpload(r *http.Request, params dto.WorkspaceParams, query, body interface{}) (string, int, error) {
	return "", http.StatusOK, nil
}
