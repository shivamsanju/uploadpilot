package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/pkg/db"
	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	"github.com/uploadpilot/uploadpilot/pkg/utils"
	webmodels "github.com/uploadpilot/uploadpilot/pkg/web/models"
)

type importHandler struct {
	impRepo db.ImportRepo
}

func NewImportHandler() *importHandler {
	return &importHandler{
		impRepo: db.NewImportRepo(),
	}
}

func (h *importHandler) GetAllImportsForUploader(w http.ResponseWriter, r *http.Request) {
	uploaderId := chi.URLParam(r, "uploaderId")
	skip, limit, search, err := utils.GetSkipLimitSearchParams(r)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	imports, totalRecords, err := h.impRepo.FindAllImportsByUploaderId(r.Context(), uploaderId, skip, limit, search)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, &webmodels.PaginatedResponse[models.Import]{
		TotalRecords: totalRecords,
		Records:      imports,
	})
}

func (h *importHandler) GetImportDetailsByID(w http.ResponseWriter, r *http.Request) {
	impID := chi.URLParam(r, "importId")
	cb, err := h.impRepo.Get(r.Context(), impID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}
