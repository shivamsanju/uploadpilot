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
	"github.com/uploadpilot/uploadpilot/web/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type workspaceHandler struct {
	wsRepo   db.WorkspaceRepo
	userRepo db.UserRepo
}

func NewWorkspaceHandler() *workspaceHandler {
	return &workspaceHandler{
		wsRepo:   db.NewWorkspaceRepo(),
		userRepo: db.NewUserRepo(),
	}
}

var ValidRoles = map[models.UserRole]bool{
	models.UserRoleOwner:       true,
	models.UserRoleContributor: true,
	models.UserRoleViewer:      true,
}

var DefaultUploaderConfig = &models.UploaderConfig{
	AllowedSources:         []models.AllowedSources{models.FileUpload},
	RequiredMetadataFields: []string{},
}

func (h *workspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	infra.Log.Info("creating workspace")
	userID := r.Header.Get("userId")
	workspace := &models.Workspace{}
	if err := render.DecodeJSON(r.Body, workspace); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	workspace.UploaderConfig = DefaultUploaderConfig
	workspace.CreatedBy = r.Header.Get("email")
	workspace.UpdatedBy = r.Header.Get("email")

	if err := validate.Struct(workspace); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

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

	addRequest := &dto.AddUserToWorkspaceRequest{}
	if err := render.DecodeJSON(r.Body, addRequest); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(addRequest); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	user, err := h.userRepo.GetUserByEmail(r.Context(), addRequest.Email)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("unknown user: %s", addRequest.Email))
		return
	}

	exists, err := h.wsRepo.CheckIfUserExistsInWorkspace(r.Context(), workspaceId, user.UserID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if exists {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("user already exists in workspace"))
		return
	}

	if !ValidRoles[addRequest.Role] {
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("unknown role: %s", addRequest.Role))
		return
	}

	workspaceUser := &models.WorkspaceUser{
		UserID: user.UserID,
		Role:   addRequest.Role,
	}

	err = h.wsRepo.AddUserToWorkspace(r.Context(), workspaceId, workspaceUser)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)

}

func (h *workspaceHandler) RemoveUserFromWorkspace(w http.ResponseWriter, r *http.Request) {
	wsId := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	userId := r.Header.Get("userId")
	err = h.wsRepo.RemoveUserFromWorkspace(r.Context(), workspaceId, userId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (h *workspaceHandler) GetAllUsersInWorkspace(w http.ResponseWriter, r *http.Request) {
	wsId := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	users, err := h.wsRepo.GetUsersInWorkspace(r.Context(), workspaceId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, users)
}

func (h *workspaceHandler) GetUploaderConfig(w http.ResponseWriter, r *http.Request) {
	wsId := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	uploaderConfig, err := h.wsRepo.GetUploaderConfig(r.Context(), workspaceId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, uploaderConfig)
}

func (h *workspaceHandler) UpdateUploaderConfig(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userId")
	wsId := chi.URLParam(r, "workspaceId")
	workspaceId, err := primitive.ObjectIDFromHex(wsId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	uploaderConfig := &models.UploaderConfig{}
	if err := render.DecodeJSON(r.Body, uploaderConfig); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(uploaderConfig); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}

	err = h.wsRepo.UpdateUploaderConfig(r.Context(), workspaceId, uploaderConfig, userID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, nil)
}

func (h *workspaceHandler) GetAllAllowedSources(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []models.AllowedSources{
		models.FileUpload,
		models.Webcamera,
		models.Audio,
		models.ScreenCapture,
		models.Box,
		models.Dropbox,
		models.Facebook,
		models.GoogleDrive,
		models.GooglePhotos,
		models.Instagram,
		models.OneDrive,
		models.Unsplash,
		models.Url,
		models.Zoom,
	})
}
