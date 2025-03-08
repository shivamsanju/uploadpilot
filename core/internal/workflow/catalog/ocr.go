package catalog

var DetectDocumentTextV1 = &ActivityMetadata{
	Name:        "DetectDocumentTextV1",
	DisplayName: "Detect Document Text V1",
	Description: `This activity extracts text content from a given document.
This activity is useful for extracting text from various document formats, such as PDF, Images, DOC, DOCX, and more.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: DetectDocumentTextV1
    saveOutput: true
`,
}

var ImageToTextV1 = &ActivityMetadata{
	Name:        "ImageToTextV1",
	DisplayName: "Image to Text V1",
	Description: `This activity extracts text content from an image.
This activity is useful for extracting text from various image formats, such as JPEG, PNG, and more.`,
	Workflow: `
- activity:
	key: uniqueKeyLessThan20Chars
	uses: ImageToTextV1
	saveOutput: true
`,
}
