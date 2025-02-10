package tasks

import "github.com/uploadpilot/uploadpilot/common/pkg/models"

type ExtractPDFContentTaskData struct {
	WorkspaceID   string `json:"workspaceID" validate:"required"`
	UploadID      string `json:"uploadID" validate:"required"`
	ProcessorID   string `json:"processorID" validate:"required"`
	IncludeImages bool   `json:"includeImages"`
}

var ExtractPDFContentTask = &Task[ExtractPDFContentTaskData]{
	Key:   models.TaskKey("ExtractPDFContent"),
	Label: "Extract PDF content",
	Description: `This task extracts text content from a given PDF file. 
	Optionally, images can be extracted if the IncludeImages flag is set. 
	The extracted content is processed and returned for further use. 
	This feature is particularly useful for text analysis, data extraction, and document processing workflows.`,
	Data: ExtractPDFContentTaskData{},
}
