package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/dto"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"github.com/uploadpilot/uploadpilot/internal/workspace"
)

type workspaceHandler struct {
	workspaceSvc *workspace.WorkspaceService
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		workspaceSvc: workspace.NewWorkspaceService(),
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

func (h *workspaceHandler) GetWorkspacesForUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserDetailsFromContext(r.Context())

	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	workspaces, err := h.workspaceSvc.GetWorkspaces(r.Context(), user.UserID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, workspaces)
}

func (h *workspaceHandler) AddUserToWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	addRequest := &dto.AddWorkspaceUser{}
	if err := render.DecodeJSON(r.Body, addRequest); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.AddUserToWorkspace(r.Context(), workspaceID, addRequest)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (h *workspaceHandler) RemoveUserFromWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	userID := chi.URLParam(r, "userId")

	err := h.workspaceSvc.RemoveUserFromWorkspace(r.Context(), workspaceID, userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *workspaceHandler) ChangeUserRoleInWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	userID := chi.URLParam(r, "userId")

	body := &dto.EditUserRole{}
	if err := render.DecodeJSON(r.Body, body); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.ChangeUserRoleInWorkspace(r.Context(), workspaceID, userID, body.Role)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, true)
}

func (h *workspaceHandler) GetAllUsersInWorkspace(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	users, err := h.workspaceSvc.GetWorkspaceUsers(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, users)
}

func (h *workspaceHandler) GetUploaderConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")
	uploaderConfig, err := h.workspaceSvc.GetUploaderConfig(r.Context(), workspaceID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, uploaderConfig)
}

func (h *workspaceHandler) UpdateUploaderConfig(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	uploaderConfig := &models.UploaderConfig{}
	if err := render.DecodeJSON(r.Body, uploaderConfig); err != nil {
		utils.HandleHttpError(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	err := h.workspaceSvc.SetUploaderConfig(r.Context(), workspaceID, uploaderConfig)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (h *workspaceHandler) GetAllAllowedSources(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, models.AllAllowedSources)
}
