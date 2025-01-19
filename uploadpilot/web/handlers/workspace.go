package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type workspaceHandler struct {
	wsRepo db.WorkspaceRepo
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		wsRepo: db.NewWorkspaceRepo(),
	}
}

func (h *workspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	infra.Log.Info("creating workspace")
	userID := r.Header.Get("userId")
	workspace := &models.Workspace{}
	if err := render.DecodeJSON(r.Body, workspace); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(workspace); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	workspace.CreatedBy = r.Header.Get("email")
	workspace.UpdatedBy = r.Header.Get("email")

	cb, err := h.wsRepo.Create(r.Context(), workspace, userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb.ID)
}

func (h *workspaceHandler) GetWorkspacesForUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	ws, err := h.wsRepo.GetWorkspacesForUser(r.Context(), userId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, ws)
}

func (h *workspaceHandler) AddUserToWorkspace(w http.ResponseWriter, r *http.Request) {
	wsId := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	workspaceUser := &models.WorkspaceUser{}
	if err := render.DecodeJSON(r.Body, workspaceUser); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(workspaceUser); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	err = h.wsRepo.AddUserToWorkspace(r.Context(), workspaceId, workspaceUser)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)

}
