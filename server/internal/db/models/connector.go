package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type S3Config struct {
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type LocalConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

type AzureConfig struct {
	AccountName string `json:"accountName"`
	AccountKey  string `json:"accountKey"`
}

type GCSConfig struct {
	ApiKey string `json:"apiKey"`
}

// StorageType defines the supported destination types
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeAzure StorageType = "azure"
	StorageTypeGCS   StorageType = "gcs"
)

type StorageConnector struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	Type        StorageType        `json:"type" validate:"required"`
	Tags        []string           `json:"tags,omitempty"`
	S3Config    *S3Config          `json:"s3Config,omitempty" `
	AzureConfig *AzureConfig       `json:"azureConfig,omitempty" `
	GCSConfig   *GCSConfig         `json:"gcsConfig,omitempty" `
	LocalConfig *LocalConfig       `json:"localConfig,omitempty" `
	CreatedAt   primitive.DateTime `json:"createdAt"`
	CreatedBy   string             `json:"createdBy"`
	UpdatedAt   primitive.DateTime `json:"updatedAt"`
	UpdatedBy   string             `json:"updatedBy"`
}
