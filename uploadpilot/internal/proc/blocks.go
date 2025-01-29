package proc

import (
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/notify"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks/pdf"
)

type TaskBlock struct {
	Key         models.TaskKey    `json:"key" validate:"required"`
	Label       string            `json:"label" validate:"required"`
	Description string            `json:"description" validate:"required"`
	TaskBuilder func() tasks.Task `json:"-"`
}

var (
	ExtractPDFContentTaskKey models.TaskKey = "extract_pdf_content"
	ExtractPDFImageTaskKey   models.TaskKey = "extract_pdf_image"
	EncodePDFContentTaskKey  models.TaskKey = "encode_pdf_content"
	WebhookTaskKey           models.TaskKey = "webhook"
	OCRTaskKey               models.TaskKey = "ocr"
)

var ProcTaskBlocks = []TaskBlock{
	{
		Key:         ExtractPDFContentTaskKey,
		Label:       "Extract PDF content",
		Description: "Extracts the PDF content from the uploaded file including images, tables, and text.",
		TaskBuilder: pdf.NewExtractPDFContentTask,
	},
	{
		Key:         EncodePDFContentTaskKey,
		Label:       "Encode PDF content",
		Description: "Encodes the PDF content into a zip file.",
		TaskBuilder: pdf.NewEncodePDFContentTask,
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
}
