package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/uploadpilot/uploadpilot/pkg/db"
	"github.com/uploadpilot/uploadpilot/pkg/db/models"
	g "github.com/uploadpilot/uploadpilot/pkg/globals"
	"github.com/uploadpilot/uploadpilot/pkg/utils"
)

type storageConnectorHandler struct {
	scRepo db.StorageConnectorRepo
}

func NewStorageConnectorHandler() *storageConnectorHandler {
	return &storageConnectorHandler{
		scRepo: db.NewStorageConnectorRepo(),
	}
}

func (h *storageConnectorHandler) GetAllStorageConnectors(w http.ResponseWriter, r *http.Request) {
	storages, err := h.scRepo.GetAll(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, storages)
}

func (h *storageConnectorHandler) GetStorageConnectorByID(w http.ResponseWriter, r *http.Request) {
	storageID := chi.URLParam(r, "id")
	cb, err := h.scRepo.Get(r.Context(), storageID)
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
	id, err := h.scRepo.Create(r.Context(), connector)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (h *storageConnectorHandler) DeleteStorageConnector(w http.ResponseWriter, r *http.Request) {
	storageID := chi.URLParam(r, "id")
	h.scRepo.Delete(r.Context(), storageID)
}
