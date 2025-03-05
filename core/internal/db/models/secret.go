package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

var keyRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

type Secret struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	WorkspaceID string    `gorm:"column:workspace_id;not null;uniqueIndex:idx_secret_workspace_id_key;type:uuid" json:"workspaceId"  validate:"required,uuid"`
	Key         string    `gorm:"column:key;not null;type:varchar(255);uniqueIndex:idx_secret_workspace_id_key" json:"key" validate:"required"`
	Value       string    `gorm:"column:value;not null" json:"value,omitempty" validate:"required"`
	Salt        string    `gorm:"column:salt;not null" json:"salt,omitempty" validate:"required"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"workspace"`
	CreatedAtColumn
	UpdatedAtColumn
	CreatedByColumn
	UpdatedByColumn
}

func (s *Secret) TableName() string {
	return "secrets"
}

func (s *Secret) BeforeCreate(tx *gorm.DB) error {
	return s.verifyEntry()
}

func (s *Secret) BeforeUpdate(tx *gorm.DB) error {
	return s.verifyEntry()
}

func (s *Secret) verifyEntry() error {
	if !keyRegex.MatchString(s.Key) {
		return errors.New("key must start with an alphabet and contain only alphanumeric characters and underscores")
	}
	s.ID = s.WorkspaceID + "_" + s.Key
	return nil
}
