package handlers

import (
	"net/http"

	"github.com/uploadpilot/core/internal/db/models"
	"github.com/uploadpilot/core/internal/dto"
	"github.com/uploadpilot/core/internal/services"
	"github.com/uploadpilot/core/pkg/utils"
)

type uploadHandler struct {
	uploadSvc    *services.UploadService
	workspaceSvc *services.WorkspaceService
}

func NewUploadHandler(uploadSvc *services.UploadService, workspaceSvc *services.WorkspaceService) *uploadHandler {
	return &uploadHandler{
		uploadSvc:    uploadSvc,
		workspaceSvc: workspaceSvc,
	}
}

func (h *uploadHandler) GetPaginatedUploads(r *http.Request, params dto.WorkspaceParams, query dto.PaginatedQuery,
	body interface{}) (*dto.PaginatedResponse[models.Upload], int, error) {

	paginationParams, err := utils.GetPaginatedQueryParams(&query)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	uploads, totalRecords, err := h.uploadSvc.GetAllUploadsForWorkspace(r.Context(), params.TenantID, params.WorkspaceID, paginationParams)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &dto.PaginatedResponse[models.Upload]{
		TotalRecords: totalRecords,
		Records:      uploads,
	}, http.StatusOK, nil
}

func (h *uploadHandler) GetUploadDetailsByID(r *http.Request, params dto.UploadParams, query interface{}, body interface{}) (*models.Upload, int, error) {
	details, err := h.uploadSvc.GetUploadDetails(r.Context(), params.TenantID, params.WorkspaceID, params.UploadID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return details, http.StatusOK, nil
}

func (h *uploadHandler) GetUploadURL(r *http.Request, params dto.UploadParams, query interface{}, body interface{}) (string, int, error) {
	url, err := h.uploadSvc.GetUploadSignedURL(r.Context(), params.TenantID, params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return url, http.StatusOK, nil
}

func (h *uploadHandler) ProcessUpload(r *http.Request, params dto.UploadParams, query interface{}, body interface{}) (string, int, error) {
	err := h.uploadSvc.ProcessUpload(r.Context(), params.TenantID, params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return "OK", http.StatusOK, nil
}

// UPLOADER API
func (h *uploadHandler) CreateUpload(r *http.Request, params dto.WorkspaceParams, query interface{}, body dto.CreateUploadRequest) (*dto.CreateUploadResponse, int, error) {
	res, err := h.uploadSvc.CreateUpload(r.Context(), params.TenantID, params.WorkspaceID, &body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return res, http.StatusOK, nil
}

func (h *uploadHandler) FinishUpload(r *http.Request, params dto.UploadParams, query interface{}, body dto.FinishUploadRequest) (bool, int, error) {
	err := h.uploadSvc.FinishUpload(r.Context(), params.TenantID, params.WorkspaceID, params.UploadID, body.Status)
	if err != nil {
		return false, http.StatusBadRequest, err
	}

	return true, http.StatusOK, nil
}
