package models

import "github.com/uploadpilot/core/internal/db/dtypes"

type Workspace struct {
	ID          string             `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string             `gorm:"column:name;not null;size:100" json:"name,omitempty"`
	Description *string            `gorm:"column:description;size:255" json:"description,omitempty"`
	Tags        dtypes.StringArray `gorm:"column:tags;type:text[]" json:"tags,omitempty"`
	TenantID    string             `gorm:"column:tenant_id;not null;type:uuid" json:"tenantId,omitempty"`
	CreatedAtColumn
	CreatedByColumn
	UpdatedAtColumn
	UpdatedByColumn
}

func (Workspace) TableName() string {
	return "workspaces"
}
