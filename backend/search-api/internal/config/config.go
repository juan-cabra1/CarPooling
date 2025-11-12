package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	Mongo       MongoConfig
	HTTP        HTTPConfig
}

type HTTPConfig struct {
	UsersAPIURL string
	TripsAPIURL string
	Timeout     int // Timeout in seconds for HTTP requests
	MaxRetries  int // Maximum number of retries for failed requests
}

type MongoConfig struct {
	URI string
	DB  string
}

func LoadConfig() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8003"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "carpooling_search"),
		},
		HTTP: HTTPConfig{
			UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
			TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
			Timeout:     getEnvInt("HTTP_TIMEOUT", 5),     // 5 seconds default
			MaxRetries:  getEnvInt("HTTP_MAX_RETRIES", 3), // 3 retries default
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Mongo.URI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}
	if c.Mongo.DB == "" {
		return fmt.Errorf("MONGO_DB is required")
	}
	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}
	if c.HTTP.UsersAPIURL == "" {
		return fmt.Errorf("USERS_API_URL is required")
	}
	if c.HTTP.TripsAPIURL == "" {
		return fmt.Errorf("TRIPS_API_URL is required")
	}
	if c.HTTP.Timeout <= 0 {
		return fmt.Errorf("HTTP_TIMEOUT must be positive")
	}
	if c.HTTP.MaxRetries < 0 {
		return fmt.Errorf("HTTP_MAX_RETRIES must be non-negative")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
