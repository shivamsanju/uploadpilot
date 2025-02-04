package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AppName            string
	WebServerPort      int
	MongoURI           string
	PostgresURI        string
	EncryptionKey      []byte
	RedisAddr          string
	RedisUsername      string
	RedisPassword      string
	DatabaseName       string
	FrontendURI        string
	AllowedOrigins     string
	RootPassword       string
	CompanionEndpoint  string
	SelfEndpoint       string
	JWTSecretKey       string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	GithubClientID     string
	GithubClientSecret string
	GithubCallbackURL  string
	TusUploadBasePath  string
	TusUploadDir       string
	S3AccessKey        string
	S3SecretKey        string
	S3BucketName       string
	S3Region           string
)

func Init() error {
	// Load the .env file
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("No .env file found, reading from environment variables instead")
	}

	portStr := os.Getenv("WEB_SERVER_PORT")
	if portStr == "" {
		portStr = "8081"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid WEB_SERVER_PORT: %w", err)
	}

	key := os.Getenv("ENCRYPTION_KEY")
	keyBytes := getValidKey(key)

	AppName = os.Getenv("APP_NAME")
	WebServerPort = port
	MongoURI = os.Getenv("MONGO_URI")
	PostgresURI = os.Getenv("POSTGRES_URI")
	EncryptionKey = keyBytes
	RedisAddr = os.Getenv("REDIS_HOST")
	RedisUsername = os.Getenv("REDIS_USER")
	RedisPassword = os.Getenv("REDIS_PASS")
	FrontendURI = os.Getenv("FRONTEND_URI")
	DatabaseName = os.Getenv("APP_NAME") + "db"
	RootPassword = os.Getenv("ROOT_PASSWORD")
	CompanionEndpoint = os.Getenv("COMPANION_ENDPOINT")
	SelfEndpoint = os.Getenv("SELF_ENDPOINT")
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

	TusUploadDir = "./tmp"
	TusUploadBasePath = SelfEndpoint + "/upload"

	return nil
}

func getValidKey(encryptionKey string) []byte {
	// Create a SHA-256 hash of the encryption key to make sure it's 32 bytes
	hash := sha256.New()
	hash.Write([]byte(encryptionKey))
	return hash.Sum(nil) // Returns a 32-byte slice (AES-256)
}
