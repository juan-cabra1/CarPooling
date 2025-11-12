package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the bookings-api microservice
// It uses environment variables with sensible defaults for development
type Config struct {
	// ServerPort is the HTTP port the server will listen on
	// Default: "8003"
	ServerPort string

	// DatabaseURL is the MySQL connection string in DSN format
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// Example: "root:password@tcp(localhost:3306)/bookings_db?charset=utf8mb4&parseTime=True&loc=Local"
	DatabaseURL string

	// JWTSecret is the secret key used for signing and validating JWT tokens
	// IMPORTANT: Must match the secret used by other microservices for authentication
	// Should be a strong, random string in production
	JWTSecret string

	// RabbitMQURL is the connection URL for RabbitMQ message broker
	// Format: amqp://user:password@host:port/
	// Example: "amqp://guest:guest@localhost:5672/"
	RabbitMQURL string

	// Environment specifies the deployment environment
	// Values: "development", "staging", "production"
	// Affects logging verbosity and error handling behavior
	Environment string

	// TripsAPIURL is the base URL for the trips-api microservice
	// Format: http://host:port (no trailing slash)
	// Example: "http://localhost:8002" or "http://trips-api:8002"
	TripsAPIURL string
}

// LoadConfig reads configuration from environment variables
// It first attempts to load from a .env file, then falls back to system environment variables
//
// Environment Variables:
//   - SERVER_PORT: HTTP server port (default: "8003")
//   - DATABASE_URL: MySQL connection DSN (required)
//   - JWT_SECRET: JWT signing secret (required in production)
//   - RABBITMQ_URL: RabbitMQ connection URL (required)
//   - TRIPS_API_URL: Base URL for trips-api (default: "http://localhost:8002")
//   - ENVIRONMENT: deployment environment (default: "development")
//
// Returns:
//   - *Config: Loaded configuration
//   - error: Error if required variables are missing or invalid
func LoadConfig() (*Config, error) {
	// Attempt to load .env file from current directory
	// This is optional - if the file doesn't exist, we continue with system env vars
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("✅ Loaded configuration from .env file")
	}

	// Build configuration from environment variables with defaults
	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8003"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		RabbitMQURL: getEnv("RABBITMQ_URL", ""),
		TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration fields are set
// Returns an error if any required field is missing or invalid
func (c *Config) Validate() error {
	// Check DatabaseURL is set
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Check RabbitMQURL is set
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}

	// Check JWTSecret is set (critical for security)
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Check TripsAPIURL is set
	if c.TripsAPIURL == "" {
		return fmt.Errorf("TRIPS_API_URL is required")
	}

	// Warn if using development environment in production
	if c.Environment == "production" && c.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("⚠️  WARNING: Using default JWT_SECRET in production is insecure!")
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv retrieves an environment variable with a fallback default value
// If the environment variable is not set, it returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
