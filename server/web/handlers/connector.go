package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/shivamsanju/uploader/internal/db/models"
	"github.com/shivamsanju/uploader/internal/db/repo"
	g "github.com/shivamsanju/uploader/pkg/globals"
	"github.com/shivamsanju/uploader/web/utils"
)

type StorageConnectorHandler interface {
	GetStorageConnectors(w http.ResponseWriter, r *http.Request)
	CreateStorageConnector(w http.ResponseWriter, r *http.Request)
}

type storageConnectorHandler struct {
}

func NewStorageConnectorHandler() StorageConnectorHandler {
	return &storageConnectorHandler{}
}

func (h *storageConnectorHandler) GetStorageConnectors(w http.ResponseWriter, r *http.Request) {
	storages, err := repo.GetStorageConnectors(r.Context())
	if err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	render.JSON(w, r, storages)
}

func (h *storageConnectorHandler) CreateStorageConnector(w http.ResponseWriter, r *http.Request) {
	g.Log.Info("creating storage")
	storage := &models.StorageConnector{}
	if err := render.DecodeJSON(r.Body, storage); err != nil {
		utils.HandleHttpError(w, r, http.StatusBadRequest, err)
		return
	}
	g.Log.Infof("adding storage: %+v", storage)
	id := repo.AddStorageConnector(r.Context(), storage)
	render.JSON(w, r, id)
}
