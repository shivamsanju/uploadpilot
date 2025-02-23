package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Common
	Port              int
	Environment       string
	CompanionEndpoint string
	UploaderEndpoint  string

	// CoreService
	CoreServiceEndpoint string
	CoreServiceAPIKey   string

	// Tus
	TusUploadBasePath   string
	TusUploadDir        string
	TusMaxFileSize      int64
	TusChunkSize        int64
	TusDisableTerminate bool
	TusDisableDownload  bool

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

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	appConfig.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}

	appConfig.Environment = os.Getenv("ENVIRONMENT")
	appConfig.CompanionEndpoint = os.Getenv("COMPANION_ENDPOINT")
	appConfig.UploaderEndpoint = os.Getenv("UPLOADER_ENDPOINT")
	appConfig.CoreServiceEndpoint = os.Getenv("CORE_SERVICE_ENDPOINT")
	appConfig.CoreServiceAPIKey = os.Getenv("CORE_SERVICE_API_KEY")
	appConfig.S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	appConfig.S3SecretKey = os.Getenv("S3_SECRET_KEY")
	appConfig.S3BucketName = os.Getenv("S3_BUCKET_NAME")
	appConfig.S3Region = os.Getenv("S3_REGION")

	// TUS Config
	maxSize := os.Getenv("TUS_MAX_FILE_SIZE")
	if maxSize == "" {
		maxSize = "1048576000" // 1GB
	}
	appConfig.TusMaxFileSize, err = strconv.ParseInt(maxSize, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TUS_MAX_FILE_SIZE: %w", err)
	}
	chunkSize := os.Getenv("TUS_CHUNK_SIZE")
	if chunkSize == "" {
		chunkSize = "10485760" // 10MB
	}
	appConfig.TusChunkSize, err = strconv.ParseInt(chunkSize, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TUS_CHUNK_SIZE: %w", err)
	}
	appConfig.TusUploadDir = os.Getenv("TUS_UPLOAD_DIR")
	if appConfig.TusUploadDir == "" {
		appConfig.TusUploadDir = "/tmp"
	}
	appConfig.TusUploadBasePath = appConfig.UploaderEndpoint + "/upload/%s"
	appConfig.TusDisableDownload = os.Getenv("TUS_DISABLE_DOWNLOAD") == "true"
	appConfig.TusDisableTerminate = os.Getenv("TUS_DISABLE_TERMINATION") == "true"

	return nil
}
