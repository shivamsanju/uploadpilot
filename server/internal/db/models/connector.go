package models

import "go.mongodb.org/mongo-driver/v2/bson"

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
	Id          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string        `json:"name"`
	Type        StorageType   `json:"type"`
	Tags        []string      `json:"tags,omitempty" `
	S3Config    *S3Config     `json:"s3Config,omitempty" `
	AzureConfig *AzureConfig  `json:"azureConfig,omitempty" `
	GCSConfig   *GCSConfig    `json:"gcsConfig,omitempty" `
	LocalConfig *LocalConfig  `json:"localConfig,omitempty" `
	CreatedAt   int64         `json:"createdAt"`
	CreatedBy   string        `json:"createdBy"`
}
