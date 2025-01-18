package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
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
}

func NewConfig() (*Config, error) {
	// Load the .env file
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("No .env file found, reading from environment variables instead")
	}

	// Parse individual environment variables
	portStr := os.Getenv("WEB_SERVER_PORT")
	if portStr == "" {
		portStr = "8081"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid WEB_SERVER_PORT: %w", err)
	}

	return &Config{
		AppName:            os.Getenv("APP_NAME"),
		WebServerPort:      port,
		MongoURI:           os.Getenv("MONGO_URI"),
		FrontendURI:        os.Getenv("FRONTEND_URI"),
		DatabaseName:       os.Getenv("APP_NAME") + "db",
		RootPassword:       os.Getenv("ROOT_PASSWORD"),
		CompanionEndpoint:  os.Getenv("COMPANION_ENDPOINT"),
		SelfEndpoint:       os.Getenv("SELF_ENDPOINT"),
		JWTSecretKey:       os.Getenv("JWT_SECRET_KEY"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubCallbackURL:  os.Getenv("GITHUB_CALLBACK_URL"),
	}, nil
}
