package models

import "github.com/shivamsanju/uploader/internal/db/models"

type DataStoreResponse struct {
	*models.DataStore
	Connector *models.StorageConnector `json:"connector"`
}
