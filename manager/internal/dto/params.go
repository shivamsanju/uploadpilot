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
