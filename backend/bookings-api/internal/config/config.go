package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	RabbitMQURL string
	Environment string
	TripsAPIURL string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		// Variables CRÍTICAS - Sin defaults, DEBEN existir (fail-fast)
		DatabaseURL: mustGetEnv("DATABASE_URL_BOOKINGS"),
		JWTSecret:   mustGetEnv("JWT_SECRET"),
		RabbitMQURL: mustGetEnv("RABBITMQ_URL"),

		// Variables NO CRÍTICAS - Con defaults razonables
		ServerPort:  getEnv("SERVER_PORT", "8003"),
		TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return cfg, nil
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
// Use only for non-critical configuration
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// mustGetEnv obtiene variable REQUERIDA o hace panic (fail-fast)
// Use for critical configuration that must be present
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("FATAL: Required environment variable " + key + " is not set")
	}
	return value
}
