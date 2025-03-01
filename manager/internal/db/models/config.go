package models

import "github.com/uploadpilot/manager/internal/db/dtypes"

type WorkspaceConfig struct {
	WorkspaceID            string             `gorm:"column:workspace_id;primaryKey;type:uuid" json:"workspaceId"`
	MaxFileSize            *int64             `gorm:"column:max_file_size" json:"maxFileSize"`
	MinFileSize            *int64             `gorm:"column:min_file_size" json:"minFileSize"`
	MaxNumberOfFiles       *int64             `gorm:"column:max_number_of_files" json:"maxNumberOfFiles"`
	MinNumberOfFiles       *int64             `gorm:"column:min_number_of_files" json:"minNumberOfFiles"`
	AllowedFileTypes       dtypes.StringArray `gorm:"column:allowed_file_types;type:text[]" json:"allowedFileTypes"`
	AllowedSources         dtypes.StringArray `gorm:"not null;type:text[];column:allowed_sources" json:"allowedSources"`
	RequiredMetadataFields dtypes.StringArray `gorm:"not null;column:required_metadata_fields;type:text[]" json:"requiredMetadataFields"`
	AllowPauseAndResume    bool               `gorm:"column:allow_pause_and_resume default:true" json:"allowPauseAndResume"`
	EnableImageEditing     bool               `gorm:"column:enable_image_editing default:false" json:"enableImageEditing"`
	UseCompression         bool               `gorm:"column:use_compression default:true" json:"useCompression"`
	UseFaultTolerantMode   bool               `gorm:"column:use_fault_tolerant_mode default:false" json:"useFaultTolerantMode"`
	AllowedOrigins         dtypes.StringArray `gorm:"column:allowed_origins;type:text[]" json:"allowedOrigins"`
	Workspace              Workspace          `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAtColumn
	UpdatedAtColumn
}

func (u *WorkspaceConfig) TableName() string {
	return "workspace_config"
}

const (
	FileUpload    string = "FileUpload"
	Audio         string = "Audio"
	Webcamera     string = "Webcamera"
	ScreenCapture string = "ScreenCapture"
	Box           string = "Box"
	Dropbox       string = "Dropbox"
	Facebook      string = "Facebook"
	GoogleDrive   string = "GoogleDrive"
	GooglePhotos  string = "GooglePhotos"
	Instagram     string = "Instagram"
	OneDrive      string = "OneDrive"
	Unsplash      string = "Unsplash"
	Url           string = "Url"
	Zoom          string = "Zoom"
)

var AllAllowedSources = dtypes.StringArray{
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

var DefaultWorkspaceConfig = &WorkspaceConfig{
	MaxFileSize:            nil,
	MinFileSize:            nil,
	MaxNumberOfFiles:       nil,
	MinNumberOfFiles:       nil,
	AllowedFileTypes:       nil,
	AllowedSources:         []string{FileUpload, Webcamera, Audio, ScreenCapture},
	RequiredMetadataFields: []string{},
	AllowPauseAndResume:    true,
	EnableImageEditing:     false,
	UseCompression:         true,
	UseFaultTolerantMode:   false,
}
