package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/pkg/db"
	"github.com/uploadpilot/uploadpilot/pkg/utils"
)

type ImportHandler interface {
	GetAllImports(w http.ResponseWriter, r *http.Request)
	GetAllImportsForUploader(w http.ResponseWriter, r *http.Request)
	GetImportDetailsByID(w http.ResponseWriter, r *http.Request)
}

type importHandler struct {
	impRepo db.ImportRepo
}

func NewImportHandler() ImportHandler {
	return &importHandler{
		impRepo: db.NewImportRepo(),
	}
}

func (h *importHandler) GetAllImports(w http.ResponseWriter, r *http.Request) {
	imports, err := h.impRepo.GetAll(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, imports)
}

func (h *importHandler) GetAllImportsForUploader(w http.ResponseWriter, r *http.Request) {
	uploaderId := chi.URLParam(r, "uploaderId")
	imports, err := h.impRepo.FindImportsByUploaderId(r.Context(), uploaderId)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, imports)
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
