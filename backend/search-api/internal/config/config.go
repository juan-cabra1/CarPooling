package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the search-api service
type Config struct {
	// Server Configuration
	ServerPort string

	// MongoDB Configuration
	MongoURI string
	MongoDB  string

	// Apache Solr Configuration
	SolrURL  string
	SolrCore string

	// Memcached Configuration
	MemcachedServers []string
	CacheTTL         int // Cache Time-To-Live in seconds

	// RabbitMQ Configuration
	RabbitMQURL string
	QueueName   string

	// External APIs
	TripsAPIURL string
	UsersAPIURL string

	// JWT Configuration
	JWTSecret string

	// Application Environment
	Environment string
}

// LoadConfig loads configuration from environment variables
// It first attempts to load from a .env file, then reads from the environment
// Returns a fully configured Config struct with defaults applied
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		// Server Configuration - default port 8004
		ServerPort: getEnv("SERVER_PORT", "8004"),

		// MongoDB Configuration
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "search_db"),

		// Apache Solr Configuration
		SolrURL:  getEnv("SOLR_URL", "http://localhost:8983/solr"),
		SolrCore: getEnv("SOLR_CORE", "trips"),

		// Memcached Configuration
		MemcachedServers: getEnvAsSlice("MEMCACHED_SERVERS", []string{"localhost:11211"}, ","),
		CacheTTL:         getEnvAsInt("CACHE_TTL", 300), // Default 5 minutes

		// RabbitMQ Configuration
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:   getEnv("QUEUE_NAME", "search.events"),

		// External APIs
		TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),

		// JWT Configuration
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Application Environment
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	// Validate required configuration fields
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate checks that all required configuration fields are set
func (c *Config) validate() error {
	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}

	if c.MongoDB == "" {
		return fmt.Errorf("MONGO_DB is required")
	}

	if c.SolrURL == "" {
		return fmt.Errorf("SOLR_URL is required")
	}

	if c.SolrCore == "" {
		return fmt.Errorf("SOLR_CORE is required")
	}

	if len(c.MemcachedServers) == 0 {
		return fmt.Errorf("MEMCACHED_SERVERS is required")
	}

	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}

	if c.QueueName == "" {
		return fmt.Errorf("QUEUE_NAME is required")
	}

	if c.TripsAPIURL == "" {
		return fmt.Errorf("TRIPS_API_URL is required")
	}

	if c.UsersAPIURL == "" {
		return fmt.Errorf("USERS_API_URL is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required for authentication")
	}

	if c.CacheTTL < 0 {
		return fmt.Errorf("CACHE_TTL must be a positive integer")
	}

	return nil
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// GetSolrFullURL returns the complete Solr URL including the core
func (c *Config) GetSolrFullURL() string {
	return fmt.Sprintf("%s/%s", c.SolrURL, c.SolrCore)
}

// Helper functions for reading environment variables

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsSlice reads an environment variable as a slice of strings or returns a default value
// The environment variable should contain comma-separated values
func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	// Split by separator and trim whitespace
	values := strings.Split(valueStr, separator)
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}

	return values
}
