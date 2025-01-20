package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AppName            string
	WebServerPort      int
	MongoURI           string
	DatabaseName       string
	FrontendURI        string
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

	AppName = os.Getenv("APP_NAME")
	WebServerPort = port
	MongoURI = os.Getenv("MONGO_URI")
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

	TusUploadDir = "./tmp"
	TusUploadBasePath = SelfEndpoint + "/upload"

	return nil
}
