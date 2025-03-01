package catalog

var WebhookV_1_0 = &ActivityMetadata{
	Name:        "Webhook@1.0",
	DisplayName: "Webhook",
	Description: `
Sends a webhook to a target URL. The webhook will be sent
with the following headers:
  - X-UploadPilot-Workspace-ID: The ID of the workspace.
  - X-UploadPilot-Processor-ID: The ID of the processor.
  - X-UploadPilot-Upload-ID: The ID of the upload.
  - X-UploadPilot-Secret: The secret provided in the webhook task data.
  `,
	Workflow: `
- activity:
  key: uniqueKeyLessThan20Chars
  input: inputActivityKey
  uses : Webhook@1.0
  with:
    url: "https://example.com/webhook"
`,
}
