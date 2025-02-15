package dto

import "github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"

type UploaderConfig struct {
	models.UploaderConfig
	ChunkSize int64 `json:"chunkSize"`
}
