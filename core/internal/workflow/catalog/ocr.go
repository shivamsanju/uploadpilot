package catalog

var DetectDocumentText = &ActivityMetadata{
	Name:        "DetectDocumentText",
	DisplayName: "Detect Document Text",
	Description: `This activity extracts text content from a given document.
This activity is useful for extracting text from various document formats, such as PDF, Images, DOC, DOCX, and more.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: DetectDocumentText
    saveOutput: true
`,
}
