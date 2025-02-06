package proc

import (
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/file"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/image"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/notify"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/pdf"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/trigger"
)

type TaskBlock struct {
	Key         models.TaskKey    `json:"key" validate:"required"`
	Label       string            `json:"label" validate:"required"`
	Description string            `json:"description" validate:"required"`
	TaskBuilder func() tasks.Task `json:"-"`
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
		TaskBuilder: trigger.NewTriggerTask,
	},
	{
		Key:         ExtractPDFContentTaskKey,
		Label:       "Extract PDF content",
		Description: "Extracts the PDF content from the uploaded file including images, tables, and text.",
		TaskBuilder: pdf.NewExtractPDFContentTask,
	},
	{
		Key:         EncryptContentTaskKey,
		Label:       "Encrypt file content",
		Description: "Encrypts the content of the uploaded file.",
		TaskBuilder: file.NewEncryptContentTask,
	},
	{
		Key:         ExtractPDFImageTaskKey,
		Label:       "Extract image from PDF",
		Description: "Extracts images from the uploaded PDF file.",
		TaskBuilder: pdf.NewExtractPDFImageTask,
	},
	{
		Key:         WebhookTaskKey,
		Label:       "Webhook",
		Description: "Executes a webhook.",
		TaskBuilder: notify.NewWebhookTask,
	},
	{
		Key:         OCRTaskKey,
		Label:       "Extract text from image",
		Description: "Extracts text from an image.",
		TaskBuilder: image.NewExtractTextFromImageTask,
	},
}
