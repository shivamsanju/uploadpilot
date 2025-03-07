package dto

import (
	"time"

	"github.com/uploadpilot/core/internal/db/dtypes"
)

type CreateProcessorRequest struct {
	Name        string             `json:"name" validate:"required,min=3,max=25,alphanumspace"`
	WorkspaceID string             `json:"workspaceId" validate:"required,uuid"`
	Triggers    dtypes.StringArray `json:"triggers" validate:"required,max=500"`
	TemplateKey string             `json:"templateKey"`
}

type EditProcRequest struct {
	Name     string             `json:"name" validate:"required,min=3,max=25,alphanumspace"`
	Triggers dtypes.StringArray `json:"triggers" validate:"required,max=500"`
}

type EnableDisableProcessorRequest struct {
	Enabled bool `json:"enabled"`
}

type WorkflowUpdate struct {
	Workflow string `json:"workflow"`
}

type TriggerWorkflowResp struct {
	WorkflowID string `json:"workflowId"`
	RunID      string `json:"runId"`
}

type WorkflowRun struct {
	ID                 string    `json:"id"`
	WorkspaceID        string    `json:"workspaceId"`
	UploadID           string    `json:"uploadId"`
	WorkflowID         string    `json:"workflowId"`
	RunID              string    `json:"runId"`
	StartTime          time.Time `json:"startTime,omitempty"`
	EndTime            time.Time `json:"endTime,omitempty"`
	WorkflowTimeMillis int64     `json:"workflowTimeMillis"`
	ExecutionTimeMilis int64     `json:"executionTimeMillis"`
	Status             string    `json:"status,omitempty"`
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

type WorkflowQuery struct {
	WorkflowID string `json:"workflowId"`
	RunID      string `json:"runId"`
}
