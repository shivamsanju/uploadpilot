package models

import "time"

type CreatedByColumn struct {
	CreatedBy string `gorm:"column:created_by;not null" json:"createdBy"`
}

type UpdatedByColumn struct {
	UpdatedBy string `gorm:"column:updated_by;not null" json:"updatedBy"`
}

type CreatedAtColumn struct {
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`
}

type UpdatedAtColumn struct {
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}
