package dto

import "github.com/uploadpilot/uploadpilot/common/pkg/types"

type CreateProcessorRequest struct {
	Name        string            `json:"name" validate:"required,max=25"`
	WorkspaceID string            `json:"workspaceId" validate:"required,uuid"`
	Triggers    types.StringArray `json:"triggers" validate:"required,min=1"`
}

type EditProcRequest struct {
	Name     string            `json:"name"`
	Triggers types.StringArray `json:"triggers"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}

type WorkflowUpdate struct {
	Workflow string `json:"workflow"`
}
