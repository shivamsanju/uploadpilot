package models

import (
	"time"

	"github.com/uploadpilot/go-core/db/pkg/dtypes"
)

type UploadLog struct {
	ID          string       `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	WorkspaceID string       `gorm:"type:uuid;column:workspace_id;not null" json:"workspaceId"`
	Data        dtypes.JSONB `gorm:"column:data;type:jsonb" json:"data"`
	Timestamp   time.Time    `gorm:"column:timestamp;type:timestamp;not null" json:"timestamp"`
	Workspace   Workspace    `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE" json:"workspace"`
}

func (UploadLog) TableName() string {
	return "upload_logs"
}
