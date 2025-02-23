package activitycatalog

type ActivityMetadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Workflow    string `json:"workflow"`
}

var ActivityCatalog = []*ActivityMetadata{
	ImageResizeV1_0,
	ImageConvertToPngV1_0,
	ImageConvertToBmpV1_0,
	ImageAddWatermarkActivityV1_0,
	ImageMetadataExtractionActivityV1_0,
	ImageBlurActivityV1_0,
	ExtractPDFContentV1_0,
	WebhookV_1_0,
}
