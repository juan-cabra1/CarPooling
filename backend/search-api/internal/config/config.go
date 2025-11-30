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
	// Intentar cargar .env desde la raíz del proyecto
	// En Docker, las variables vienen del docker-compose, así que esto falla silenciosamente
	_ = godotenv.Load()

	cfg := &Config{
		// Variables CRÍTICAS - Sin defaults, DEBEN existir (fail-fast)
		Mongo: MongoConfig{
			URI: mustGetEnv("MONGO_URI_SEARCH"),
			DB:  getEnv("MONGO_DB_SEARCH", "carpooling_search"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:       mustGetEnv("RABBITMQ_URL"),
			QueueName: getEnv("QUEUE_NAME", "search.events"),
		},
		JWT: JWTConfig{
			Secret: mustGetEnv("JWT_SECRET"),
		},

		// Variables NO CRÍTICAS - Con defaults razonables
		ServerPort: getEnv("SERVER_PORT", "8004"),
		Solr: SolrConfig{
			URL:  getEnv("SOLR_URL", "http://localhost:8983/solr"),
			Core: getEnv("SOLR_CORE", "carpooling_trips"),
		},
		Memcached: MemcachedConfig{
			Servers: getEnvSlice("MEMCACHED_SERVERS", []string{"localhost:11211"}),
		},
		HTTP: HTTPConfig{
			UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
			TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
			Timeout:     getEnvInt("HTTP_TIMEOUT", 5),     // 5 seconds default
			MaxRetries:  getEnvInt("HTTP_MAX_RETRIES", 3), // 3 retries default
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// mustGetEnv obtiene variable REQUERIDA o hace panic (fail-fast)
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("FATAL: Required environment variable " + key + " is not set")
	}
	return value
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
