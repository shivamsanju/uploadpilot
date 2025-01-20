package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/utils"
	"github.com/uploadpilot/uploadpilot/web/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type importHandler struct {
	impRepo db.ImportRepo
}

func NewImportHandler() *importHandler {
	return &importHandler{
		impRepo: db.NewImportRepo(),
	}
}

func (h *importHandler) GetAllImportsForWorkspace(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceId")
	workspaceID, err := primitive.ObjectIDFromHex(wID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	skip, limit, search, err := utils.GetSkipLimitSearchParams(r)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	imports, totalRecords, err := h.impRepo.FindAllImportsForWorkspace(r.Context(), workspaceID, skip, limit, search)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, &dto.PaginatedResponse[models.Import]{
		TotalRecords: totalRecords,
		Records:      imports,
	})
}

func (h *importHandler) GetImportDetailsByID(w http.ResponseWriter, r *http.Request) {
	impID := chi.URLParam(r, "importId")
	importID, err := primitive.ObjectIDFromHex(impID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	cb, err := h.impRepo.Get(r.Context(), importID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}
