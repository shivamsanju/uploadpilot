package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/manager/internal/dto"
	"github.com/uploadpilot/uploadpilot/manager/internal/svc"
	"github.com/uploadpilot/uploadpilot/manager/internal/utils"
)

type uploadHandler struct {
	upSvc svc.UploadService
}

func NewUploadHandler() *uploadHandler {
	return &uploadHandler{
		upSvc: *svc.NewUploadService(),
	}
}

func (h *uploadHandler) GetPaginatedUploads(w http.ResponseWriter, r *http.Request) {
	workspaceID := chi.URLParam(r, "workspaceId")

	skip, limit, search, err := utils.GetSkipLimitSearchParams(r)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	uploads, totalRecords, err := h.upSvc.GetAllUploads(r.Context(), workspaceID, skip, limit, search)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, &dto.PaginatedResponse[models.Upload]{
		TotalRecords: totalRecords,
		Records:      uploads,
	})
}

func (h *uploadHandler) GetUploadDetailsByID(w http.ResponseWriter, r *http.Request) {
	uploadID := chi.URLParam(r, "uploadId")
	workspaceID := chi.URLParam(r, "workspaceId")

	details, err := h.upSvc.GetUploadDetails(r.Context(), workspaceID, uploadID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, details)
}

func (h *uploadHandler) GetUploadLogs(w http.ResponseWriter, r *http.Request) {
	uploadID := chi.URLParam(r, "uploadId")

	logs, err := h.upSvc.GetLogs(r.Context(), uploadID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	render.JSON(w, r, logs)
}
