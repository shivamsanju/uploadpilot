package models

type AllowedSources string

const (
	FileUpload    AllowedSources = "FileUpload"
	Audio         AllowedSources = "Audio"
	Webcamera     AllowedSources = "Webcamera"
	ScreenCapture AllowedSources = "ScreenCapture"
	Box           AllowedSources = "Box"
	Dropbox       AllowedSources = "Dropbox"
	Facebook      AllowedSources = "Facebook"
	GoogleDrive   AllowedSources = "GoogleDrive"
	GooglePhotos  AllowedSources = "GooglePhotos"
	Instagram     AllowedSources = "Instagram"
	OneDrive      AllowedSources = "OneDrive"
	Unsplash      AllowedSources = "Unsplash"
	Url           AllowedSources = "Url"
	Zoom          AllowedSources = "Zoom"
)

type UploaderConfig struct {
	MaxFileSize            int              `bson:"maxFileSize" json:"maxFileSize"`
	MinFileSize            int              `bson:"minFileSize" json:"minFileSize"`
	MaxNumberOfFiles       int              `bson:"maxNumberOfFiles" json:"maxNumberOfFiles"`
	MinNumberOfFiles       int              `bson:"minNumberOfFiles" json:"minNumberOfFiles"`
	MaxTotalFileSize       int              `bson:"maxTotalFileSize" json:"maxTotalFileSize"`
	AllowedFileTypes       []string         `bson:"allowedFileTypes" json:"allowedFileTypes"`
	AllowedSources         []AllowedSources `bson:"allowedSources" json:"allowedSources" validate:"required"`
	RequiredMetadataFields []string         `bson:"requiredMetadataFields" json:"requiredMetadataFields" validate:"required"`
	AllowPauseAndResume    bool             `bson:"allowPauseAndResume" json:"allowPauseAndResume"`
	EnableImageEditing     bool             `bson:"enableImageEditing" json:"enableImageEditing"`
	UseCompression         bool             `bson:"useCompression" json:"useCompression"`
	UseFaultTolerantMode   bool             `bson:"useFaultTolerantMode" json:"useFaultTolerantMode"`
	AuthEndpoint           string           `bson:"authEndpoint" json:"authEndpoint"`
}

var AllAllowedSources = []AllowedSources{
	FileUpload,
	Webcamera,
	Audio,
	ScreenCapture,
	Box,
	Dropbox,
	Facebook,
	GoogleDrive,
	GooglePhotos,
	Instagram,
	OneDrive,
	Unsplash,
	Url,
	Zoom,
}
