package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	"github.com/shivamsanju/uploader/internal/web/utils"
	g "github.com/shivamsanju/uploader/pkg/globals"
)

type storageConnectorHandler struct {
	scRepo repo.StorageConnectorRepo
}

func NewStorageConnectorHandler() *storageConnectorHandler {
	return &storageConnectorHandler{
		scRepo: repo.NewStorageConnectorRepo(),
	}
}

func (h *storageConnectorHandler) GetAllStorageConnectors(w http.ResponseWriter, r *http.Request) {
	storages, err := h.scRepo.GetStorageConnectors(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, storages)
}

func (h *storageConnectorHandler) GetStorageConnectorByID(w http.ResponseWriter, r *http.Request) {
	storageID := chi.URLParam(r, "id")
	cb, err := h.scRepo.GetStorageConnector(r.Context(), storageID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, cb)
}

func (h *storageConnectorHandler) CreateStorageConnector(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating storage")
	connector := &models.StorageConnector{}
	if err := render.DecodeJSON(r.Body, connector); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(connector); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}
	g.Log.Infof("adding storage connector: %+v", connector)
	connector.CreatedBy = r.Header.Get("email")
	connector.UpdatedBy = r.Header.Get("email")
	id, err := h.scRepo.CreateStorageConnector(r.Context(), connector)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (h *storageConnectorHandler) DeleteStorageConnector(w http.ResponseWriter, r *http.Request) {
	storageID := chi.URLParam(r, "id")
	h.scRepo.DeleteStorageConnector(r.Context(), storageID)
}
