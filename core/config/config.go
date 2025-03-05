package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	AppName     string `mapstructure:"APP_NAME"`

	// Webserver
	Port           int      `mapstructure:"PORT"`
	SelfEndpoint   string   `mapstructure:"SELF_ENDPOINT"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
	FrontendURI    string   `mapstructure:"FRONTEND_URI"`

	// Auth
	SuperTokensEndpoint string `mapstructure:"SUPERTOKENS_ENDPOINT"`
	SupertokensAPIKey   string `mapstructure:"SUPERTOKENS_API_KEY"`
	GoogleClientID      string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret  string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleCallbackURL   string `mapstructure:"GOOGLE_CALLBACK_URL"`
	GithubClientID      string `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret  string `mapstructure:"GITHUB_CLIENT_SECRET"`
	GithubCallbackURL   string `mapstructure:"GITHUB_CALLBACK_URL"`
	APIBaseAuthPath     string `mapstructure:"API_BASE_AUTH_PATH"`
	WebsiteBaseAuthPath string `mapstructure:"WEBSITE_BASE_AUTH_PATH"`

	// Api key
	ApiKeyEncryptionKey  string `mapstructure:"API_KEY_ENCRYPTION_KEY"`
	ApiKeyEncryptionSalt string `mapstructure:"API_KEY_ENCRYPTION_SALT"`

	// Encryption
	EncryptionKey  string `mapstructure:"ENCRYPTION_KEY"`
	EncryptionSalt string `mapstructure:"ENCRYPTION_SALT"`

	// Database
	PostgresURI string `mapstructure:"POSTGRES_URI"`
	Database    string `mapstructure:"DATABASE"`
	UseCache    bool   `mapstructure:"USE_CACHE"`

	// Redis
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisTLS      bool   `mapstructure:"REDIS_TLS"`

	// Storage S3
	S3AccessKey  string `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey  string `mapstructure:"S3_SECRET_KEY"`
	S3BucketName string `mapstructure:"S3_BUCKET_NAME"`
	S3Region     string `mapstructure:"S3_REGION"`

	// Temporal
	TemporalNamespace string `mapstructure:"TEMPORAL_NAMESPACE"`
	TemporalHostPort  string `mapstructure:"TEMPORAL_HOST_PORT"`
	TemporalAPIKey    string `mapstructure:"TEMPORAL_API_KEY"`
}

var AppConfig *Config

func LoadConfig(path, filename, ext string) error {
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("APP_NAME", "UploadPilot")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("SELF_ENDPOINT", "http://localhost:8080")
	viper.SetDefault("ALLOWED_ORIGINS", []string{"http://localhost:3000"})
	viper.SetDefault("FRONTEND_URI", "http://localhost:3000")
	viper.SetDefault("SUPERTOKENS_ENDPOINT", "https://try.supertokens.com")
	viper.SetDefault("API_BASE_AUTH_PATH", "/auth")
	viper.SetDefault("WEBSITE_BASE_AUTH_PATH", "/auth")
	viper.SetDefault("API_KEY_ENCRYPTION_KEY", "thisisaverylooooooooooongsecret")
	viper.SetDefault("API_KEY_ENCRYPTION_SALT", "thisisaverylooooooooooongsalt")
	viper.SetDefault("ENCRYPTION_KEY", "thisisaverylooooooooooongsecret")
	viper.SetDefault("ENCRYPTION_SALT", "thisisaverylooooooooooongsalt")
	viper.SetDefault("POSTGRES_URI", "localhost:5432")
	viper.SetDefault("DATABASE", "uploadpilotdb")
	viper.SetDefault("USE_CACHE", false)
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_TLS", false)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal into Config struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return nil
}
