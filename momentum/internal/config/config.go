package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/uploadpilot/uploadpilot/common/pkg/kms"
)

var (
	Port          int
	PostgresURI   string
	SecretKey     string
	EncryptionKey []byte

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

	key := os.Getenv("ENCRYPTION_KEY")
	keyBytes, err := kms.GetValidKey(key)
	if err != nil {
		return fmt.Errorf("invalid ENCRYPTION_KEY: %w", err)
	}
	EncryptionKey = keyBytes

	// Temporal
	TemporalNamespace = os.Getenv("TEMPORAL_NAMESPACE")
	TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	TemporalAPIKey = os.Getenv("TEMPORAL_API_KEY")

	return nil
}
