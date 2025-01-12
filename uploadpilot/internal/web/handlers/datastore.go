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

type datastoreHandler struct {
	dsRepo repo.DataStoreRepo
	scRepo repo.StorageConnectorRepo
}

func NewDatastoreHandler() *datastoreHandler {
	return &datastoreHandler{
		dsRepo: repo.NewDataStoreRepo(),
		scRepo: repo.NewStorageConnectorRepo(),
	}
}

func (h *datastoreHandler) GetAllDatastores(w http.ResponseWriter, r *http.Request) {
	datastores, err := h.dsRepo.GetDataStores(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, datastores)
}

func (h *datastoreHandler) GetDatastoreByID(w http.ResponseWriter, r *http.Request) {
	datastoreID := chi.URLParam(r, "id")
	dataStore, err := h.dsRepo.GetDataStore(r.Context(), datastoreID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}

	g.Log.Infof("found datastore: %+v", &dataStore)
	render.JSON(w, r, dataStore)
}

func (h *datastoreHandler) CreateDatastore(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating datastore")
	datastore := &models.DataStore{}
	if err := render.DecodeJSON(r.Body, datastore); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(datastore); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		utils.HandleHttpError(w, r, http.StatusBadRequest, fmt.Errorf("validation error: %v", errors))
		return
	}
	g.Log.Infof("creating datastore: %+v", datastore)
	datastore.CreatedBy = r.Header.Get("email")
	datastore.UpdatedBy = r.Header.Get("email")
	id, err := h.dsRepo.CreateDataStore(r.Context(), datastore)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, id)
}

func (h *datastoreHandler) DeleteDatastore(w http.ResponseWriter, r *http.Request) {
	datastoreID := chi.URLParam(r, "id")
	h.dsRepo.DeleteDataStore(r.Context(), datastoreID)
}
