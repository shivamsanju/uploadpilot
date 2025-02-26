package models

import (
	"time"

	"github.com/uploadpilot/go-core/db/pkg/dtypes"
)

type Upload struct {
	ID             string       `gorm:"column:id;primaryKey;default:uuid_generate_v4();type:uuid" json:"id"`
	WorkspaceID    string       `gorm:"column:workspace_id;type:uuid;not null" json:"workspaceId"`
	Status         UploadStatus `gorm:"column:status;not null" json:"status"`
	StatusReason   string       `gorm:"column:status_reason" json:"statusReason"`
	Metadata       dtypes.JSONB `gorm:"column:metadata;type:jsonb" json:"metadata"`
	FileName       string       `gorm:"column:file_name;not null" json:"fileName"`
	FileType       string       `gorm:"column:file_type;not null" json:"fileType"`
	StoredFileName string       `gorm:"column:stored_file_name;not null" json:"storedFileName"`
	Size           int64        `gorm:"column:size;not null" json:"size"`
	URL            string       `gorm:"column:url" json:"url"`
	ProcessedURL   string       `gorm:"column:processed_url" json:"processedUrl"`
	StartedAt      time.Time    `gorm:"column:started_at;default:now()" json:"startedAt"`
	FinishedAt     time.Time    `gorm:"column:finished_at" json:"finishedAt"`
	Workspace      Workspace    `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE" json:"workspace"`
}

type UploadStatus string

const (
	UploadStatusStarted             UploadStatus = "Started"
	UploadStatusSkipped             UploadStatus = "Skipped"
	UploadStatusInProgress          UploadStatus = "In Progress"
	UploadStatusComplete            UploadStatus = "Uploaded"
	UploadStatusFailed              UploadStatus = "Failed"
	UploadStatusCancelled           UploadStatus = "Cancelled"
	UploadStatusProcessing          UploadStatus = "Processing"
	UploadStatusProcessingFailed    UploadStatus = "Processing Failed"
	UploadStatusProcessingComplete  UploadStatus = "Processing Complete"
	UploadStatusProcessingCancelled UploadStatus = "Processing Cancelled"
	UploadStatusDeleted             UploadStatus = "Deleted"
)

var UploadTerminalStates = []UploadStatus{
	UploadStatusSkipped,
	UploadStatusComplete,
	UploadStatusFailed,
	UploadStatusCancelled,
	UploadStatusProcessingFailed,
	UploadStatusProcessingComplete,
	UploadStatusDeleted,
}

var UploadNonTerminalStates = []UploadStatus{
	UploadStatusStarted,
	UploadStatusInProgress,
	UploadStatusProcessing,
}

var UploadAllStates = append(UploadTerminalStates, UploadNonTerminalStates...)
