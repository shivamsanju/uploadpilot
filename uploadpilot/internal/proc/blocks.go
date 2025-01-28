package proc

import (
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
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
	EncodePDFContentTaskKey  models.TaskKey = "encode_pdf_content"
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
}
