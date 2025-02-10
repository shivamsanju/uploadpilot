package dto

import "github.com/uploadpilot/uploadpilot/common/pkg/types"

type EditProcRequest struct {
	Name     string            `json:"name"`
	Triggers types.StringArray `json:"triggers"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}
