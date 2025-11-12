package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	Mongo       MongoConfig
	UsersAPIURL string
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
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
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
	if c.UsersAPIURL == "" {
		return fmt.Errorf("USERS_API_URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
