package models

import (
	"time"

	"github.com/uploadpilot/core/internal/db/dtypes"
)

type Upload struct {
	ID            string       `gorm:"column:id;primaryKey;default:uuid_generate_v4();type:uuid" json:"id"`
	WorkspaceID   string       `gorm:"column:workspace_id;type:uuid;not null" json:"workspaceId,omitempty"`
	FileName      string       `gorm:"column:file_name" json:"fileName,omitempty"`
	ContentType   string       `gorm:"column:content_type" json:"contentType,omitempty"`
	ContentLength int64        `gorm:"column:content_length" json:"contentLength,omitempty"`
	Metadata      dtypes.JSONB `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	Status        UploadStatus `gorm:"column:status;not null" json:"status,omitempty"`
	StartedAt     time.Time    `gorm:"column:started_at;default:now()" json:"startedAt,omitempty"`
	FinishedAt    time.Time    `gorm:"column:finished_at" json:"finishedAt,omitempty"`
	Workspace     Workspace    `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE" json:"-"`
}

type UploadStatus string

const (
	UploadStatusCreated             UploadStatus = "Created"
	UploadStatusSkipped             UploadStatus = "Skipped"
	UploadStatusFinished            UploadStatus = "Finished"
	UploadStatusFailed              UploadStatus = "Failed"
	UploadStatusCancelled           UploadStatus = "Cancelled"
	UploadStatusProcessing          UploadStatus = "Processing"
	UploadStatusProcessingFailed    UploadStatus = "Processing Failed"
	UploadStatusProcessingComplete  UploadStatus = "Processing Complete"
	UploadStatusProcessingCancelled UploadStatus = "Processing Cancelled"
	UploadStatusDeleted             UploadStatus = "Deleted"
	UploadStatusTimedOut            UploadStatus = "Timed Out"
)

var UploadTerminalStates = []UploadStatus{
	UploadStatusSkipped,
	UploadStatusFinished,
	UploadStatusFailed,
	UploadStatusCancelled,
	UploadStatusProcessingFailed,
	UploadStatusProcessingComplete,
	UploadStatusDeleted,
	UploadStatusTimedOut,
}

var UploadNonTerminalStates = []UploadStatus{
	UploadStatusCreated,
	UploadStatusProcessing,
}

var UploadAllStates = append(UploadTerminalStates, UploadNonTerminalStates...)
