package models

import (
	"time"

	"gorm.io/gorm"
)

type CreatedByColumn struct {
	CreatedBy string `gorm:"column:created_by;not null" json:"createdBy,omitempty"`
}

type UpdatedByColumn struct {
	UpdatedBy string `gorm:"column:updated_by;not null" json:"updatedBy,omitempty"`
}

type CreatedAtColumn struct {
	CreatedAt time.Time `gorm:"column:created_at;not null,default:now()" json:"createdAt,omitempty"`
}

type UpdatedAtColumn struct {
	UpdatedAt time.Time `gorm:"column:updated_at;not null,default:now()" json:"updatedAt,omitempty"`
}

func (t *UpdatedAtColumn) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

func (t *CreatedAtColumn) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	return nil
}
