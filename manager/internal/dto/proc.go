package dto

import "github.com/uploadpilot/uploadpilot/manager/internal/db/types"

type EditProcRequest struct {
	Name     string            `json:"name"`
	Triggers types.StringArray `json:"triggers"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}

type UpdateProcTaskRequest struct {
	Canvas types.JSONB          `json:"canvas"`
	Data   types.EncryptedJSONB `json:"data"`
}
