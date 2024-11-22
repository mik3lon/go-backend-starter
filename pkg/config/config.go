package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration variables.
type Config struct {
	AddressPort string
	DatabaseDSN string

	GoogleClientId     string
	PrivateKeyPEM      string
	PrivateKeyPassword string
	PublicKeyPEM       string

	S3Region      string
	S3Endpoint    string
	S3ImageBucket string
	AppEnv        string
}

// LoadConfig loads environment variables from a .env file and populates the Config struct.
func LoadConfig() *Config {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return load()
}

// LoadTestConfig loads environment variables from a .env file and populates the Config struct.
func LoadTestConfig() *Config {
	// Load the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return load()
}

func load() *Config {
	return &Config{
		AddressPort:        getEnv("ADDRESS_PORT", ":8081"),                                      // Default to :8081 if not set
		DatabaseDSN:        getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/db"), // Default DSN
		GoogleClientId:     getEnv("GOOGLE_CLIENT_ID", ""),
		PrivateKeyPEM:      getEnv("USER_PRIVATE_PEM_FILE", ""),
		PrivateKeyPassword: getEnv("USER_PRIVATE_PEM_PASSWORD", ""),
		PublicKeyPEM:       getEnv("USER_PUBLIC_PEM_FILE", ""),
		S3Endpoint:         getEnv("AWS_S3_ENDPOINT", ""),
		S3Region:           getEnv("AWS_S3_REGION", "us-east-1"),
		S3ImageBucket:      getEnv("AWS_S3_IMAGE_BUCKET", ""),
		AppEnv:             getEnv("APP_ENV", "test"),
	}
}

// getEnv gets an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
