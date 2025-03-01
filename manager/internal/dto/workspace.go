package dto

import (
	"github.com/uploadpilot/manager/internal/db/models"
)

type WorkspaceConfig struct {
	models.WorkspaceConfig
	ChunkSize int64 `json:"chunkSize"`
}
