package models

import "github.com/uploadpilot/uploadpilot/internal/db/types"

type Workspace struct {
	ID          string            `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string            `gorm:"column:name;not null;size:100" json:"name"`
	Description *string           `gorm:"column:description;size:255" json:"description"`
	Tags        types.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	At
	By
}

func (Workspace) TableName() string {
	return "workspaces"
}

type UserWorkspace struct {
	UserID      string    `gorm:"column:user_id;type:uuid;primaryKey;index" json:"userId"`
	WorkspaceID string    `gorm:"column:workspace_id;type:uuid;primaryKey;index" json:"workspaceId"`
	Role        UserRole  `gorm:"column:role;not null" json:"role"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
	At
}

func (UserWorkspace) TableName() string {
	return "user_workspaces"
}
