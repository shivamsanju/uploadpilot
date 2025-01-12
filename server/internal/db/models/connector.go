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
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Type        StorageType        `bson:"type" json:"type" validate:"required"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	S3Config    *S3Config          `bson:"s3Config,omitempty" json:"s3Config,omitempty"`
	AzureConfig *AzureConfig       `bson:"azureConfig,omitempty" json:"azureConfig,omitempty"`
	GCSConfig   *GCSConfig         `bson:"gcsConfig,omitempty" json:"gcsConfig,omitempty"`
	LocalConfig *LocalConfig       `bson:"localConfig,omitempty" json:"localConfig,omitempty"`
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
	CreatedBy   string             `bson:"createdBy" json:"createdBy"`
	UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   string             `bson:"updatedBy" json:"updatedBy"`
}
