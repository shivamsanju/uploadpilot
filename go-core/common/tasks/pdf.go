package tasks

var ExtractPDFContentTask = &Task{
	Name:        "ExtractPDFContent",
	DisplayName: "Extract PDF content",
	Description: `This task extracts text content from a given PDF file. 
	Optionally, images can be extracted if the IncludeImages flag is set. 
	The extracted content is processed and returned for further use. 
	This feature is particularly useful for text analysis, data extraction, and document processing workflows.`,
}
