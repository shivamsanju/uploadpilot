package dto

import "time"

type CreateUploadRequest struct {
	ID        string                 `json:"id,omitempty"`
	FileName  string                 `json:"fileName" validate:"required"`
	FileType  string                 `json:"fileType,omitempty"`
	Size      int64                  `json:"size" validate:"required"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Status    string                 `json:"status,omitempty"`
	StartedAt time.Time              `json:"startedAt" validate:"required"`
}

type FinishUploadRequest struct {
	FinishedAt time.Time `json:"finishedAt" validate:"required"`
	Status     string    `json:"status,omitempty" validate:"required"`
	FileType   string    `json:"fileType,omitempty"`
	Size       int64     `json:"size,omitempty"`
}
