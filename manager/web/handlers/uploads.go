package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/phuslu/log"
	commonutils "github.com/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/manager/internal/svc/upload"
	"github.com/uploadpilot/manager/internal/utils"
)

type uploadHandler struct {
	uploadSvc upload.Service
}

func NewUploadHandler(uploadSvc *upload.Service) *uploadHandler {
	return &uploadHandler{
		uploadSvc: *uploadSvc,
	}
}

func (h *uploadHandler) GetPaginatedUploads(
	ctx context.Context,
	params dto.WorkspaceParams,
	query dto.PaginatedQuery,
	body interface{},
) (*dto.PaginatedResponse[models.Upload], int, error) {

	paginationParams, err := utils.GetPaginatedQueryParams(&query)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	uploads, totalRecords, err := h.uploadSvc.GetAllUploads(ctx, params.WorkspaceID, paginationParams)
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

func (h *uploadHandler) CreateUpload(
	ctx context.Context,
	params dto.WorkspaceParams,
	query interface{},
	body dto.CreateUploadRequest,
) (string, int, error) {
	log.Info().Msgf("#sss upload: %+v", body)

	var upload models.Upload
	if err := commonutils.ConvertDTOToModel(&body, &upload); err != nil {
		return "", http.StatusUnprocessableEntity, err
	}
	if err := h.uploadSvc.CreateUpload(ctx, params.WorkspaceID, &upload); err != nil {
		return "", http.StatusBadRequest, err
	}

	return upload.ID, http.StatusOK, nil
}

func (h *uploadHandler) FinishUpload(
	ctx context.Context,
	params dto.UploadParams,
	query interface{},
	body dto.FinishUploadRequest,
) (bool, int, error) {
	err := h.uploadSvc.FinishUpload(ctx, params.WorkspaceID, params.UploadID, &body)
	if err != nil {
		return false, http.StatusBadRequest, err
	}

	return true, http.StatusOK, nil
}

func (h *uploadHandler) GetUploadURL(
	ctx context.Context,
	params dto.UploadParams,
	query interface{},
	body interface{},
) (string, int, error) {
	url, err := h.uploadSvc.GetUploadSignedURL(ctx, params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	return url, http.StatusOK, nil
}

func (h *uploadHandler) ProcessUpload(
	ctx context.Context,
	params dto.UploadParams,
	query interface{},
	body interface{},
) (string, int, error) {
	err := h.uploadSvc.ProcessUpload(ctx, params.WorkspaceID, params.UploadID)
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	return "OK", http.StatusOK, nil
}
