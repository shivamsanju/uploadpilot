package models

import "github.com/uploadpilot/uploadpilot/common/pkg/types"

type Processor struct {
	ID          string            `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string            `gorm:"column:name;not null" json:"name"`
	WorkspaceID string            `gorm:"column:workspace_id;not null;type:uuid" json:"workspaceId"`
	Triggers    types.StringArray `gorm:"column:triggers;not null;type:text[]" json:"triggers"`
	Workflow    string            `gorm:"column:workflow;type:text;not null; default:''" json:"workflow"`
	Enabled     bool              `gorm:"column:enabled;not null;default:true" json:"enabled"`
	Workspace   Workspace         `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
	At
	By
}

func (*Processor) TableName() string {
	return "processors"
}
