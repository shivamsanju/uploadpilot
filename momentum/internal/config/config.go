package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string

	// Temporal
	TemporalNamespace string
	TemporalHostPort  string
	TemporalAPIKey    string

	// S3
	S3AccessKey  string
	S3SecretKey  string
	S3BucketName string
	S3Region     string
}

var appConfig *Config

func GetAppConfig() *Config {
	if appConfig == nil {
		panic("please call BuildConfig first")
	}
	return appConfig
}

func BuildConfig() error {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("No .env file found, reading from environment variables instead")
	}

	appConfig = &Config{}
	appConfig.Environment = os.Getenv("ENVIRONMENT")

	appConfig.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	appConfig.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	appConfig.S3BucketName = os.Getenv("S3_BUCKET_NAME")
	appConfig.S3Region = os.Getenv("S3_REGION")

	appConfig.TemporalNamespace = os.Getenv("TEMPORAL_NAMESPACE")
	appConfig.TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	appConfig.TemporalAPIKey = os.Getenv("TEMPORAL_API_KEY")

	return nil
}
