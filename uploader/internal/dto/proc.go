package dto

type TriggerworkflowReq struct {
	Workflow string `json:"workflow"`
}

type TriggerWorkflowResp struct {
	WorkflowID string `json:"workflowId"`
	RunID      string `json:"runId"`
}
