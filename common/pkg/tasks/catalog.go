package tasks

import (
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
)

type TaskBlock struct {
	Key         models.TaskKey `json:"key" validate:"required"`
	Label       string         `json:"label" validate:"required"`
	Description string         `json:"description" validate:"required"`
}

var (
	TriggerTaskKey           models.TaskKey = "trigger"
	ExtractPDFContentTaskKey models.TaskKey = "extract_pdf_content"
	ExtractPDFImageTaskKey   models.TaskKey = "extract_pdf_image"
	EncryptContentTaskKey    models.TaskKey = "encrypt_content"
	WebhookTaskKey           models.TaskKey = "webhook"
	OCRTaskKey               models.TaskKey = "ocr"
)

var ProcTaskBlocks = []TaskBlock{
	{
		Key:         TriggerTaskKey,
		Label:       "",
		Description: "",
	},
	{
		Key:         ExtractPDFContentTaskKey,
		Label:       "Extract PDF content",
		Description: "Extracts the PDF content from the uploaded file including images, tables, and text.",
	},
	{
		Key:         EncryptContentTaskKey,
		Label:       "Encrypt file content",
		Description: "Encrypts the content of the uploaded file.",
	},
	{
		Key:         ExtractPDFImageTaskKey,
		Label:       "Extract image from PDF",
		Description: "Extracts images from the uploaded PDF file.",
	},
	{
		Key:         WebhookTaskKey,
		Label:       "Webhook",
		Description: "Executes a webhook.",
	},
	{
		Key:         OCRTaskKey,
		Label:       "Extract text from image",
		Description: "Extracts text from an image.",
	},
}
