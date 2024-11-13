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
}

// LoadConfig loads environment variables from a .env file and populates the Config struct.
func LoadConfig() *Config {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		AddressPort: getEnv("ADDRESS_PORT", ":8081"),                                      // Default to :8081 if not set
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/db"), // Default DSN
	}
}

// getEnv gets an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
