package dto

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type WorkflowTask struct {
	TaskKey         models.TaskKey   `json:"taskKey"`
	EncryptedData   string           `json:"encryptedData"`
	Retry           uint             `json:"retry"`
	TimeoutMs       uint64           `json:"timeoutMs"`
	ContinueOnError bool             `json:"continueOnError"`
	DependsOn       []models.TaskKey `json:"dependsOn"`
}

type TriggerworkflowReq struct {
	Tasks []WorkflowTask `json:"tasks"`
}

type TriggerworkflowResp struct {
	WorkflowID string `json:"workflowId"`
}
