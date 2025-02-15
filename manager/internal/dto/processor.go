package dto

type CreateProcessorRequest struct {
	Name        string   `json:"name" validate:"required,max=25"`
	WorkspaceID string   `json:"workspaceId" validate:"required,uuid"`
	Triggers    []string `json:"triggers" validate:"required,min=1"`
}

type EditProcRequest struct {
	Name     string   `json:"name"`
	Triggers []string `json:"triggers"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}

type WorkflowUpdate struct {
	Workflow string `json:"workflow"`
}
