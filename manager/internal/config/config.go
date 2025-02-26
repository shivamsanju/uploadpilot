package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AppName               string
	Environment           string
	Port                  int
	PostgresURI           string
	RedisAddr             string
	RedisUsername         string
	RedisPassword         string
	RedisTLS              bool
	DatabaseName          string
	FrontendURI           string
	AllowedOrigins        string
	JWTSecretKey          string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleCallbackURL     string
	GithubClientID        string
	GithubClientSecret    string
	GithubCallbackURL     string
	S3AccessKey           string
	S3SecretKey           string
	S3BucketName          string
	S3Region              string
	ApiKeyEncryptionKey   string
	ApiKeyEncryptionSalt  string
	SecretsEncryptionKey  string
	SecretsEncryptionSalt string

	// Temporal
	TemporalNamespace string
	TemporalHostPort  string
	TemporalAPIKey    string
)

func Init() error {
	// Load the .env file
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

	uPortStr := os.Getenv("UPLOADER_SERVER_PORT")
	if uPortStr == "" {
		uPortStr = "8081"
	}

	AppName = os.Getenv("APP_NAME")
	Environment = os.Getenv("ENVIRONMENT")
	PostgresURI = os.Getenv("POSTGRES_URI")
	RedisAddr = os.Getenv("REDIS_HOST")
	RedisUsername = os.Getenv("REDIS_USER")
	RedisPassword = os.Getenv("REDIS_PASS")
	RedisTLS = os.Getenv("REDIS_TLS") == "true"
	FrontendURI = os.Getenv("FRONTEND_URI")
	DatabaseName = os.Getenv("APP_NAME") + "db"
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	GoogleCallbackURL = os.Getenv("GOOGLE_CALLBACK_URL")
	GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	GithubCallbackURL = os.Getenv("GITHUB_CALLBACK_URL")
	S3AccessKey = os.Getenv("S3_ACCESS_KEY")
	S3SecretKey = os.Getenv("S3_SECRET_KEY")
	S3BucketName = os.Getenv("S3_BUCKET_NAME")
	S3Region = os.Getenv("S3_REGION")
	AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	ApiKeyEncryptionKey = os.Getenv("API_KEY_ENCRYPTION_KEY")
	ApiKeyEncryptionSalt = os.Getenv("API_KEY_ENCRYPTION_SALT")
	SecretsEncryptionKey = os.Getenv("SECRETS_ENCRYPTION_KEY")
	SecretsEncryptionSalt = os.Getenv("SECRETS_ENCRYPTION_SALT")

	// Temporal
	TemporalNamespace = os.Getenv("TEMPORAL_NAMESPACE")
	TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	TemporalAPIKey = os.Getenv("TEMPORAL_API_KEY")

	return nil
}
