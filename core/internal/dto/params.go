package dto

type TenantParams struct {
	TenantID string `json:"tenantId" validate:"required,uuid"`
}

type WorkspaceParams struct {
	TenantID    string `json:"tenantId" validate:"required,uuid"`
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
}

type ProcessorParams struct {
	TenantID    string `json:"tenantId" validate:"required,uuid"`
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
	ProcessorID string `json:"processorId" validate:"required,uuid"`
}

type WorkflowRunParams struct {
	TenantID    string `json:"tenantId" validate:"required,uuid"`
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
	ProcessorID string `json:"processorId" validate:"required,uuid"`
	WorkflowID  string `json:"workflowId" validate:"required,uuid"`
	RunID       string `json:"runId" validate:"required,uuid"`
}

type UploadParams struct {
	TenantID    string `json:"tenantId" validate:"required,uuid"`
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
	UploadID    string `json:"uploadId" validate:"required,uuid"`
}

type ApiKeyParams struct {
	TenantID string `json:"tenantId" validate:"required,uuid"`
	ApiKeyID string `json:"apiKeyId" validate:"required,uuid"`
}

type PaginatedQuery struct {
	Offset              string `json:"offset" validate:"omitempty,integer"`
	Limit               string `json:"limit" validate:"omitempty,integer"`
	Search              string `json:"search,omitempty" validate:"omitempty,max=100"`
	CaseSensitiveSearch string `json:"caseSensitiveSearch,omitempty"`
	Filter              string `json:"filter,omitempty" validate:"omitempty,keyvaluepairs,max=300"`
	Sort                string `json:"sort,omitempty" validate:"omitempty,sort"`
}
