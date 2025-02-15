package models

import "github.com/uploadpilot/uploadpilot/go-core/db/pkg/dtypes"

type UploaderConfig struct {
	WorkspaceID            string             `gorm:"column:workspace_id;primaryKey;type:uuid" json:"workspaceId"`
	MaxFileSize            *int               `gorm:"column:max_file_size" json:"maxFileSize"`
	MinFileSize            *int               `gorm:"column:min_file_size" json:"minFileSize"`
	MaxNumberOfFiles       *int               `gorm:"column:max_number_of_files" json:"maxNumberOfFiles"`
	MinNumberOfFiles       *int               `gorm:"column:min_number_of_files" json:"minNumberOfFiles"`
	AllowedFileTypes       dtypes.StringArray `gorm:"column:allowed_file_types;type:text[]" json:"allowedFileTypes"`
	AllowedSources         dtypes.StringArray `gorm:"not null;type:text[];column:allowed_sources" json:"allowedSources"`
	RequiredMetadataFields dtypes.StringArray `gorm:"not null;column:required_metadata_fields;type:text[]" json:"requiredMetadataFields"`
	AllowPauseAndResume    bool               `gorm:"column:allow_pause_and_resume default:true" json:"allowPauseAndResume"`
	EnableImageEditing     bool               `gorm:"column:enable_image_editing default:false" json:"enableImageEditing"`
	UseCompression         bool               `gorm:"column:use_compression default:true" json:"useCompression"`
	UseFaultTolerantMode   bool               `gorm:"column:use_fault_tolerant_mode default:false" json:"useFaultTolerantMode"`
	AuthEndpoint           *string            `gorm:"column:auth_endpoint" json:"authEndpoint"`
	Workspace              Workspace          `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
	At
}

func (u *UploaderConfig) TableName() string {
	return "uploader_configs"
}

type AllowedSources string

func (s AllowedSources) String() string {
	return string(s)
}

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
