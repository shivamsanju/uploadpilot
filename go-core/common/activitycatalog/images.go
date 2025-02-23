package activitycatalog

var ImageResizeV1_0 = &ActivityMetadata{
	Name:        "ImageResize@v1.0",
	DisplayName: "Resize Image",
	Description: `This Activity resizes an image to the specified dimensions while maintaining aspect ratio. 
	Users can define width and height parameters, and an optional quality setting. 
	This is essential for optimizing images for web usage, thumbnails, and performance improvements.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageResize@v1.0
    with:
      width: 100
      height: 100
`,
}

var ImageConvertToPngV1_0 = &ActivityMetadata{
	Name:        "ImageConvertToPng@v1.0",
	DisplayName: "Convert Image to PNG",
	Description: `This Activity converts an image to PNG format. 
	Compression options are available for lossy formats. 
	This Activity is useful for optimizing storage and compatibility across platforms.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageConvertToPng@v1.0
`,
}

var ImageConvertToJpegV1_0 = &ActivityMetadata{
	Name:        "ImageConvertToJpeg@v1.0",
	DisplayName: "Convert Image to JPEG",
	Description: `This Activity converts an image to JPEG format. 
	Compression options are available for lossy formats. 
	This Activity is useful for optimizing storage and compatibility across platforms.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageConvertToJpeg@v1.0
`,
}

var ImageConvertToBmpV1_0 = &ActivityMetadata{
	Name:        "ImageConvertToBmp@v1.0",
	DisplayName: "Convert Image to BMP",
	Description: `This Activity converts an image to BMP format. 
	Compression options are available for lossy formats. 
	This Activity is useful for optimizing storage and compatibility across platforms.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageConvertToBmp@v1.0
`,
}

var ImageAddWatermarkActivityV1_0 = &ActivityMetadata{
	Name:        "ImageAddWatermark@v1.0",
	DisplayName: "Apply Watermark",
	Description: `This Activity applies a watermark (text or image) to an image. 
	Positioning, opacity, and scaling options are available. 
	This helps in branding, copyright protection, and content ownership verification.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageAddWatermark@v1.0
    with:
      text: "Watermark"
      opacity: 0.5
	  size: 40
`,
}

var ImageMetadataExtractionActivityV1_0 = &ActivityMetadata{
	Name:        "ImageMetadataExtraction@v1.0",
	DisplayName: "Extract Image Metadata",
	Description: `This Activity extracts metadata (EXIF, IPTC, XMP) from an image file. 
	Includes details like camera settings, GPS location, and timestamps. 
	Useful for digital forensics, asset management, and content analysis.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageMetadataExtraction@v1.0
`,
}

var ImageBlurActivityV1_0 = &ActivityMetadata{
	Name:        "ImageBlur@v1.0",
	DisplayName: "Blur Image",
	Description: `This Activity applies a Gaussian blur to an image. 
	Users can define the blur radius and quality settings. 
	This is useful for optimizing images for web usage, thumbnails, and performance improvements.`,
	Workflow: `
- activity:
    key: uniqueKeyLessThan20Chars
    uses: ImageBlur@v1.0
    with:
      radius: 2
`,
}
