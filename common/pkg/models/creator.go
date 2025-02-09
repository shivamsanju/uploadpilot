package models

import "time"

type At struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

type By struct {
	CreatedBy string `gorm:"column:created_by;not null" json:"createdBy"`
	UpdatedBy string `gorm:"column:updated_by;not null" json:"updatedBy"`
}
