package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
}

type WorkspaceUser struct {
	UserID string   `bson:"userId" json:"userId" validate:"required"`
	Role   UserRole `bson:"role" json:"role" validate:"required"`
}

type Workspace struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Name           string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Tags           []string           `bson:"tags" json:"tags"`
	Users          []WorkspaceUser    `bson:"users" json:"users"`
	UploaderConfig *UploaderConfig    `bson:"uploaderConfig" json:"uploaderConfig" validate:"required"`
	CreatedAt      primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy      string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy      string             `bson:"updatedBy" json:"updatedBy"`
}

type WorkspaceUserWithDetails struct {
	UserID string   `bson:"userId" json:"userId" validate:"required"`
	Role   UserRole `bson:"role" json:"role" validate:"required"`
	Name   string   `bson:"name" json:"name"`
	Email  string   `bson:"email" json:"email"`
}
