package models

import "github.com/uploadpilot/uploadpilot/common/pkg/types"

type TaskKey string

type Task struct {
	ID              string               `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProcessorID     string               `gorm:"column:processor_id;not null;type:uuid" json:"processorId"`
	Key             TaskKey              `gorm:"column:key;not null" json:"key"`
	Label           string               `gorm:"column:label;not null" json:"label"`
	Name            string               `gorm:"column:name;not null" json:"name"`
	Position        int                  `gorm:"column:position;not null" json:"position"`
	Data            types.EncryptedJSONB `gorm:"column:data;type:text;not null" json:"data"`
	Retries         uint                 `gorm:"column:retries;not null;default:0" json:"retries"`
	TimeoutMs       uint64               `gorm:"column:timeout_ms;not null;default:0" json:"timeoutMs"`
	ContinueOnError bool                 `gorm:"column:continue_on_error;not null;default:false" json:"continueOnError"`
	Enabled         bool                 `gorm:"column:enabled;not null;default:true" json:"enabled"`
	Processor       Processor            `gorm:"foreignKey:ProcessorID;constraint:OnDelete:CASCADE" json:"processor"`
}
