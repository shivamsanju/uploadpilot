package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	webmodels "github.com/shivamsanju/uploader/internal/web/models"
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
	cb, err := h.dsRepo.GetDataStore(r.Context(), datastoreID)
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	connector, err := h.scRepo.GetStorageConnector(r.Context(), cb.ConnectorID.Hex())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	g.Log.Infof("found datastore: %+v", &connector)
	dsResp := webmodels.DataStoreResponse{
		DataStore: cb,
		Connector: connector,
	}
	render.JSON(w, r, dsResp)
}

func (h *datastoreHandler) CreateDatastore(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating datastore")
	datastore := &models.DataStore{}
	if err := render.DecodeJSON(r.Body, datastore); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
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
