package dto

import "github.com/uploadpilot/uploadpilot/common/pkg/models"

type UploaderConfig struct {
	models.UploaderConfig
	ChunkSize int64 `json:"chunkSize"`
}
