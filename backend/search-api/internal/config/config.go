package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	Mongo      MongoConfig
	Solr       SolrConfig
	Memcached  MemcachedConfig
	RabbitMQ   RabbitMQConfig
	HTTP       HTTPConfig
	JWT        JWTConfig
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

type SolrConfig struct {
	URL  string
	Core string
}

type MemcachedConfig struct {
	Servers []string
}

type RabbitMQConfig struct {
	URL       string
	QueueName string
}

type JWTConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	// Load .env file (ignore error if not exists)
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8004"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "carpooling_search"),
		},
		Solr: SolrConfig{
			URL:  getEnv("SOLR_URL", "http://localhost:8983/solr"),
			Core: getEnv("SOLR_CORE", "carpooling_trips"),
		},
		Memcached: MemcachedConfig{
			Servers: getEnvSlice("MEMCACHED_SERVERS", []string{"localhost:11211"}),
		},
		RabbitMQ: RabbitMQConfig{
			URL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			QueueName: getEnv("QUEUE_NAME", "search.events"),
		},
		HTTP: HTTPConfig{
			UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
			TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
			Timeout:     getEnvInt("HTTP_TIMEOUT", 5),     // 5 seconds default
			MaxRetries:  getEnvInt("HTTP_MAX_RETRIES", 3), // 3 retries default
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "dev-secret-key"),
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
	if c.Solr.URL == "" {
		return fmt.Errorf("SOLR_URL is required")
	}
	if c.Solr.Core == "" {
		return fmt.Errorf("SOLR_CORE is required")
	}
	if len(c.Memcached.Servers) == 0 {
		return fmt.Errorf("MEMCACHED_SERVERS is required")
	}
	if c.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.RabbitMQ.QueueName == "" {
		return fmt.Errorf("QUEUE_NAME is required")
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
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
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

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma for multiple servers
		var result []string
		for _, v := range splitAndTrim(value, ",") {
			if v != "" {
				result = append(result, v)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for _, c := range s {
		if string(c) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
