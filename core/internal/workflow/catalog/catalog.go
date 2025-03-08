package catalog

type ActivityMetadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Workflow    string `json:"workflow"`
}

var ActivityCatalog = []*ActivityMetadata{
	ImageConverterPNGV1,
	DetectDocumentTextV1,
	ImageToTextV1,
	HTTP_V_01,
	ImageConvertToBmpV1_0,
	ImageAddWatermarkActivityV1_0,
	ImageMetadataExtractionActivityV1_0,
	ImageBlurActivityV1_0,
	ImageResizeV1_0,
}
