package dto

import (
	"github.com/uploadpilot/core/internal/db/models"
)

type WorkspaceConfig struct {
	models.WorkspaceConfig
	ChunkSize int64 `json:"chunkSize"`
}
