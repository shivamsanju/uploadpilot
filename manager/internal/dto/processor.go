package dto

import (
	"time"

	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/dtypes"
)

type CreateProcessorRequest struct {
	Name        string             `json:"name" validate:"required,max=25"`
	WorkspaceID string             `json:"workspaceId" validate:"required,uuid"`
	Triggers    dtypes.StringArray `json:"triggers"`
	TemplateKey string             `json:"templateKey"`
}

type EditProcRequest struct {
	Name     string             `json:"name"`
	Triggers dtypes.StringArray `json:"triggers"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}

type WorkflowUpdate struct {
	Workflow string `json:"workflow"`
}

type WorkflowRun struct {
	ID              string    `json:"id"`
	WorkflowID      string    `json:"workflowId"`
	RunID           string    `json:"runId"`
	StartTime       time.Time `json:"startTime,omitempty"`
	EndTime         time.Time `json:"endTime,omitempty"`
	DurationSeconds int64     `json:"durationSeconds,omitempty"`
	Status          string    `json:"status,omitempty"`
}

type WorkflowRunLogs struct {
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"eventType"`
	Details   string    `json:"details"`
}

type ProcessorTemplate struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}
