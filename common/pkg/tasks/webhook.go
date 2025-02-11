package tasks

type WebhookTaskData struct {
	WorkspaceID string `json:"workspaceID" validate:"required"`
	UploadID    string `json:"uploadID" validate:"required"`
	ProcessorID string `json:"processorID" validate:"required"`
	Url         string `json:"url" validate:"required"`
	Secret      string `json:"secret"`
}

var WebhookTask = &Task[WebhookTaskData]{
	Key:   "Webhook",
	Label: "Webhook",
	Description: `
Sends a webhook to a target URL. The webhook will be sent
with the following headers:
  - X-UploadPilot-Workspace-ID: The ID of the workspace.
  - X-UploadPilot-Processor-ID: The ID of the processor.
  - X-UploadPilot-Upload-ID: The ID of the upload.
  - X-UploadPilot-Secret: The secret provided in the webhook task data.
  `,
	Data: WebhookTaskData{},
}
