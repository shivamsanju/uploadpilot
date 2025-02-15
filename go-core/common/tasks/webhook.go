package tasks

var WebhookTask = &Task{
	Name:        "Webhook",
	DisplayName: "Webhook",
	Description: `
Sends a webhook to a target URL. The webhook will be sent
with the following headers:
  - X-UploadPilot-Workspace-ID: The ID of the workspace.
  - X-UploadPilot-Processor-ID: The ID of the processor.
  - X-UploadPilot-Upload-ID: The ID of the upload.
  - X-UploadPilot-Secret: The secret provided in the webhook task data.
  `,
}
