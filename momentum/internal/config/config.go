package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port          int
	PostgresURI   string
	SecretKey     string
	EncryptionKey string

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
		portStr = "8083"
	}

	Port, err = strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid PORT: %w", err)
	}

	PostgresURI = os.Getenv("POSTGRES_URI")
	SecretKey = os.Getenv("SECRET_KEY")

	EncryptionKey = os.Getenv("ENCRYPTION_KEY")

	// Temporal
	TemporalNamespace = os.Getenv("TEMPORAL_NAMESPACE")
	TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	TemporalAPIKey = os.Getenv("TEMPORAL_API_KEY")

	return nil
}
