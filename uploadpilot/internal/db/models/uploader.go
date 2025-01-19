package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	AllowedFileTypes       []string         `bson:"allowedFileTypes" json:"allowedFileTypes" validate:"required"`
	AllowedSources         []AllowedSources `bson:"allowedSources" json:"allowedSources" validate:"required"`
	RequiredMetadataFields []string         `bson:"requiredMetadataFields" json:"requiredMetadataFields"`
	AllowPauseAndResume    bool             `bson:"allowPauseAndResume" json:"allowPauseAndResume"`
	EnableImageEditing     bool             `bson:"enableImageEditing" json:"enableImageEditing"`
	UseCompression         bool             `bson:"useCompression" json:"useCompression"`
	UseFaultTolerantMode   bool             `bson:"useFaultTolerantMode" json:"useFaultTolerantMode"`
}

// DataStore is the model for data stores
// Connector is duplicated here to improve reads on the cost of double writes
// But that is a fair trade as there will not be many connectors
type DataStore struct {
	Name          string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	ConnectorID   primitive.ObjectID `bson:"connectorId" json:"connectorId" validate:"required"`
	ConnectorName string             `bson:"connectorName" json:"connectorName"`
	ConnectorType string             `bson:"connectorType" json:"connectorType"`
	Bucket        string             `bson:"bucket" json:"bucket" validate:"required"`
}

type Uploader struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	WorkspaceID primitive.ObjectID `bson:"workspaceId" json:"workspaceId" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Description string             `bson:"description" json:"description" validate:"max=500"`
	Tags        []string           `bson:"tags" json:"tags"`
	Config      *UploaderConfig    `bson:"config" json:"config" validate:"required"`
	DataStore   *DataStore         `bson:"dataStore" json:"dataStore" validate:"required"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
}
