package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/manager/internal/db/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/upload"
	"github.com/uploadpilot/manager/internal/svc/workspace"
	"github.com/uploadpilot/manager/internal/utils"
)

type uploadHandler struct {
	uploadSvc    *upload.Service
	workspaceSvc *workspace.Service
}

func NewUploadHandler(uploadSvc *upload.Service, workspaceSvc *workspace.Service) *uploadHandler {
	return &uploadHandler{
		uploadSvc:    uploadSvc,
		workspaceSvc: workspaceSvc,
	}
}

func (h *uploadHandler) GetPaginatedUploads(
	r *http.Request,
	params dto.WorkspaceParams,
	query dto.PaginatedQuery,
	body interface{},
) (*dto.PaginatedResponse[models.Upload], int, error) {

	paginationParams, err := utils.GetPaginatedQueryParams(&query)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	uploads, totalRecords, err := h.uploadSvc.GetAllUploads(r.Context(), params.WorkspaceID, paginationParams)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &dto.PaginatedResponse[models.Upload]{
		TotalRecords: totalRecords,
		Records:      uploads,
	}, http.StatusOK, nil
}

func (h *uploadHandler) GetUploadDetailsByID(w http.ResponseWriter, r *http.Request) {
	uploadID := chi.URLParam(r, "uploadId")
	workspaceID := chi.URLParam(r, "workspaceId")

	details, err := h.uploadSvc.GetUploadDetails(r.Context(), workspaceID, uploadID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, details)
}

func (h *uploadHandler) GetUploadURL(
	r *http.Request,
	params dto.UploadParams,
	query interface{},
	body interface{},
) (string, int, error) {
	url, err := h.uploadSvc.GetUploadSignedURL(r.Context(), params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return url, http.StatusOK, nil
}

func (h *uploadHandler) ProcessUpload(
	r *http.Request,
	params dto.UploadParams,
	query interface{},
	body interface{},
) (string, int, error) {
	statusCode, err := h.verifySubscription(r.Context(), params.WorkspaceID)
	if err != nil {
		return "", statusCode, err
	}
	err = h.uploadSvc.ProcessUpload(r.Context(), params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return "OK", http.StatusOK, nil
}

// UPLOADER API
func (h *uploadHandler) CreateUpload(
	r *http.Request,
	params dto.WorkspaceParams,
	query interface{},
	body dto.CreateUploadRequest,
) (string, int, error) {
	statusCode, err := h.verifySubscription(r.Context(), params.WorkspaceID)
	if err != nil {
		return "", statusCode, err
	}
	var upload models.Upload
	if err := utils.ConvertDTOToModel(&body, &upload); err != nil {
		return "", http.StatusUnprocessableEntity, err
	}
	if err := h.uploadSvc.CreateUpload(r.Context(), params.WorkspaceID, &upload); err != nil {
		return "", http.StatusBadRequest, err
	}

	return upload.ID, http.StatusOK, nil
}

func (h *uploadHandler) FinishUpload(
	r *http.Request,
	params dto.UploadParams,
	query interface{},
	body dto.FinishUploadRequest,
) (bool, int, error) {
	statusCode, err := h.verifySubscription(r.Context(), params.WorkspaceID)
	if err != nil {
		return false, statusCode, err
	}
	err = h.uploadSvc.FinishUpload(r.Context(), params.WorkspaceID, params.UploadID, &body)
	if err != nil {
		return false, http.StatusBadRequest, err
	}

	return true, http.StatusOK, nil
}

func (h *uploadHandler) verifySubscription(ctx context.Context, workspaceID string) (int, error) {
	return http.StatusOK, nil
}
