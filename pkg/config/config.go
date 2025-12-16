package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all our application settings
type Config struct {
	LinkedInEmail    string
	LinkedInPassword string
	HeadlessMode     bool
}

// LoadConfig reads the .env file and returns a Config struct
func LoadConfig() *Config {
	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// 2. Return the populated struct
	return &Config{
		LinkedInEmail:    getEnv("LINKEDIN_EMAIL", ""),
		LinkedInPassword: getEnv("LINKEDIN_PASSWORD", ""),
		HeadlessMode:     getEnv("HEADLESS_MODE", "false") == "true",
	}
}

// Helper function to get a string from env or return default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
