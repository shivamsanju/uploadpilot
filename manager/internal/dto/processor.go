package dto

import (
	"time"

	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/dtypes"
)

type CreateProcessorRequest struct {
	Name        string             `json:"name" validate:"required,max=25"`
	WorkspaceID string             `json:"workspaceId" validate:"required,uuid"`
	Triggers    dtypes.StringArray `json:"triggers" validate:"required,min=1"`
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
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	DurationSeconds int64     `json:"durationSeconds"`
	Status          string    `json:"status"`
}

type WorkflowRunLogs struct {
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"eventType"`
	Details   string    `json:"details"`
}
