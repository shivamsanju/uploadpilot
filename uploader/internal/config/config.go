package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
)

var (
	// Common
	Port              int
	PostgresURI       string
	EncryptionKey     string
	RedisAddr         string
	RedisUsername     string
	RedisPassword     string
	RedisTLS          bool
	CompanionEndpoint string
	UploaderEndpoint  string

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

	//Event Bus
	EventBusRedisConfig *pubsub.RedisConfig
)

func Init() error {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("No .env file found, reading from environment variables instead")
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	Port, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}

	EncryptionKey = os.Getenv("ENCRYPTION_KEY")
	PostgresURI = os.Getenv("POSTGRES_URI")
	RedisAddr = os.Getenv("REDIS_HOST")
	RedisUsername = os.Getenv("REDIS_USER")
	RedisPassword = os.Getenv("REDIS_PASS")
	RedisTLS = os.Getenv("REDIS_TLS") == "true"
	CompanionEndpoint = os.Getenv("COMPANION_ENDPOINT")
	UploaderEndpoint = os.Getenv("UPLOADER_ENDPOINT")
	S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	S3SecretKey = os.Getenv("S3_SECRET_KEY")
	S3BucketName = os.Getenv("S3_BUCKET_NAME")
	S3Region = os.Getenv("S3_REGION")

	// TUS Config
	maxSize := os.Getenv("TUS_MAX_FILE_SIZE")
	if maxSize == "" {
		maxSize = "1048576000" // 1GB
	}
	TusMaxFileSize, err = strconv.ParseInt(maxSize, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TUS_MAX_FILE_SIZE: %w", err)
	}
	chunkSize := os.Getenv("TUS_CHUNK_SIZE")
	if chunkSize == "" {
		chunkSize = "10485760" // 10MB
	}
	TusChunkSize, err = strconv.ParseInt(chunkSize, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid TUS_CHUNK_SIZE: %w", err)
	}
	TusUploadDir = os.Getenv("TUS_UPLOAD_DIR")
	if TusUploadDir == "" {
		TusUploadDir = "/tmp"
	}
	TusUploadBasePath = UploaderEndpoint + "/upload"
	TusDisableDownload = os.Getenv("TUS_DISABLE_DOWNLOAD") == "true"
	TusDisableTerminate = os.Getenv("TUS_DISABLE_TERMINATION") == "true"

	// Event Bus
	ebRedisAddr := os.Getenv("EVENT_BUS_REDIS_ADDR")
	ebRedisUsername := os.Getenv("EVENT_BUS_REDIS_USERNAME")
	ebRedisPassword := os.Getenv("EVENT_BUS_REDIS_PASSWORD")
	ebRedisTls := os.Getenv("EVENT_BUS_REDIS_TLS") == "true"

	EventBusRedisConfig = &pubsub.RedisConfig{
		Addr:     &ebRedisAddr,
		Username: &ebRedisUsername,
		Password: &ebRedisPassword,
		TLS:      ebRedisTls,
	}

	return nil
}
