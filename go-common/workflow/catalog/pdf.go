package catalog

var ExtractPDFContentV1_0 = &ActivityMetadata{
	Name:        "ExtractPDFContent@v1.0",
	DisplayName: "Extract PDF content",
	Description: `This activity extracts text content from a given PDF file. 
	Optionally, images can be extracted if the IncludeImages flag is set. 
	The extracted content is processed and returned for further use. 
	This feature is particularly useful for text analysis, data extraction, and document processing workflows.`,
	Workflow: `
- activity:
  key: uniqueKeyLessThan20Chars
  uses: ExtractPDFContent@v1.0
  with:
    includeImages: true
`,
}
