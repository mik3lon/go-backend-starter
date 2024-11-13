package test

import (
	"bytes"
	"github.com/joho/godotenv"
	"go-boilerplate/internal/pkg/infrastructure/kernel"
	"go-boilerplate/pkg/config"
	"log"
	"net/http/httptest"
	"os"
)

var k *kernel.Kernel

func LoadTestEnv() {
	// Load the .env file
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cnf := &config.Config{
		AddressPort: getEnv("ADDRESS_PORT", ":8081"),                                      // Default to :8081 if not set
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://user:password@localhost:5432/db"), // Default DSN
	}

	k = kernel.Init(cnf)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func DoHttpRequest(method string, target string, body *bytes.Buffer) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)

	httpRecorder := httptest.NewRecorder()

	k.Router.Handler().ServeHTTP(httpRecorder, r)

	return httpRecorder
}
