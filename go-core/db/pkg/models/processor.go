package models

import "github.com/uploadpilot/go-core/db/pkg/dtypes"

type Processor struct {
	ID                        string             `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name                      string             `gorm:"column:name;not null" json:"name,omitempty"`
	WorkspaceID               string             `gorm:"column:workspace_id;not null;type:uuid" json:"workspaceId,omitempty"`
	Triggers                  dtypes.StringArray `gorm:"column:triggers;not null;type:text[]" json:"triggers,omitempty"`
	Workflow                  string             `gorm:"column:workflow;type:text;not null; default:''" json:"workflow,omitempty"`
	MaxRetries                int32              `gorm:"column:max_retries;not null;default:3" json:"maxRetries,omitempty"`
	RetryInitialIntervalS     uint64             `gorm:"column:retry_initial_interval_s;not null;default:1" json:"retryInitialIntervalS,omitempty"`
	RetryBackoffCoefficient   float64            `gorm:"column:retry_backoff_coefficient;not null;default:2.0" json:"retryBackoffCoefficient,omitempty"`
	RetryMaxIntervalS         uint64             `gorm:"column:retry_max_interval_s;not null;default:60" json:"retryMaxIntervalS,omitempty"`
	WorkflowExecutionTimeoutS uint64             `gorm:"column:workflow_execution_timeout_s;not null;default:3600" json:"workflowExecutionTimeoutS,omitempty"`
	WorkflowRunTimeoutS       uint64             `gorm:"column:workflow_run_timeout_s;not null;default:3600" json:"workflowRunTimeoutS,omitempty"`
	TaskRunTimeoutS           uint64             `gorm:"column:task_run_timeout_s;not null;default:600" json:"taskRunTimeoutS,omitempty"`
	Enabled                   bool               `gorm:"column:enabled;not null;default:true" json:"enabled,omitempty"`
	Workspace                 Workspace          `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAtColumn
	UpdatedAtColumn
	CreatedByColumn
	UpdatedByColumn
}

func (*Processor) TableName() string {
	return "processors"
}
