package models

import (
	"time"
)

type UploadLog struct {
	ID          string         `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	WorkspaceID string         `gorm:"type:uuid;column:workspace_id;not null" json:"workspaceId"`
	UploadID    string         `gorm:"type:uuid;column:upload_id;not null" json:"uploadId"`
	ProcessorID *string        `gorm:"type:uuid;column:processor_id" json:"processorId"`
	TaskID      *string        `gorm:"column:task_id;type:varchar(255)" json:"taskId"`
	Level       UploadLogLevel `gorm:"column:level;type:upload_log_level;not null" json:"level"`
	Timestamp   time.Time      `gorm:"column:timestamp;type:timestamp;not null" json:"timestamp"`
	Message     string         `gorm:"column:message;type:text;not null" json:"message"`
	Workspace   Workspace      `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
	Upload      Upload         `gorm:"foreignKey:UploadID;constraint:OnDelete:CASCADE" json:"upload"`
}

func (UploadLog) TableName() string {
	return "upload_logs"
}

type UploadLogLevel string

const (
	UploadLogLevelInfo  UploadLogLevel = "info"
	UploadLogLevelWarn  UploadLogLevel = "warn"
	UploadLogLevelError UploadLogLevel = "error"
)
