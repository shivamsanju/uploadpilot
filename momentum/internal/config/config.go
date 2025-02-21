package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Environment       string
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

	Environment = os.Getenv("ENVIRONMENT")
	// Temporal
	TemporalNamespace = os.Getenv("TEMPORAL_NAMESPACE")
	TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	TemporalAPIKey = os.Getenv("TEMPORAL_API_KEY")

	return nil
}
