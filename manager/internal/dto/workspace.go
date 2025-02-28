package dto

import (
	"github.com/uploadpilot/go-core/db/pkg/models"
)

type WorkspaceConfig struct {
	models.WorkspaceConfig
	ChunkSize int64 `json:"chunkSize"`
}
