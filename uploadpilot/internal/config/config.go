package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	WebServerPort int
	MongoURI      string
	DatabaseName  string
	FrontendURI   string
	RootPassword  string
	CompanionURI  string
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
		AppName:       os.Getenv("APP_NAME"),
		WebServerPort: port,
		MongoURI:      os.Getenv("MONGO_URI"),
		FrontendURI:   os.Getenv("FRONTEND_URI"),
		DatabaseName:  os.Getenv("APP_NAME") + "DB",
		RootPassword:  os.Getenv("ROOT_PASSWORD"),
		CompanionURI:  os.Getenv("COMPANION_URI"),
	}, nil
}
