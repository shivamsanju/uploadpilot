package dto

type WorkspaceParams struct {
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
}

type ProcessorParams struct {
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
	ProcessorID string `json:"processorId" validate:"required,uuid"`
}

type UploadParams struct {
	WorkspaceID string `json:"workspaceId" validate:"required,uuid"`
	UploadID    string `json:"uploadId" validate:"required,uuid"`
}

type PaginatedQuery struct {
	Offset              string `json:"offset" validate:"omitempty,integer"`
	Limit               string `json:"limit" validate:"omitempty,integer"`
	Search              string `json:"search,omitempty" validate:"omitempty"`
	CaseSensitiveSearch string `json:"caseSensitiveSearch,omitempty"`
	Filter              string `json:"filter,omitempty" validate:"omitempty,keyvaluepairs"`
	Sort                string `json:"sort,omitempty" validate:"omitempty,sort"`
}
