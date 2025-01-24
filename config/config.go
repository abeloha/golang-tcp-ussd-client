package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort string
	Username   string
	Password   string
	ClientID   string
}

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig() Config {
	// Try to load .env file
	if err := loadEnvFile(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	return Config{
		ServerHost: getEnvWithDefault("SERVER_HOST", "localhost"),
		ServerPort: getEnvWithDefault("SERVER_PORT", "8080"),
		Username:   getEnvWithDefault("USERNAME", "admin"),
		Password:   getEnvWithDefault("PASSWORD", "password"),
		ClientID:   getEnvWithDefault("CLIENT_ID", "client123"),
	}
}

// loadEnvFile attempts to load the .env file from the current directory
// or parent directories
func loadEnvFile() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Try to find .env file in current or parent directories
	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			return godotenv.Load(envFile)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return godotenv.Load() // try loading from default locations
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
