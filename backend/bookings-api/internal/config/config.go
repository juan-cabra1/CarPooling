package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort string

	// JWT
	JWTSecret string

	// External Services
	TripsAPIURL string
	UsersAPIURL string

	// RabbitMQ
	RabbitMQURL      string
	RabbitMQExchange string
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "carpooling_reservations"),

		// Server
		ServerPort: getEnv("SERVER_PORT", "8003"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// External Services
		TripsAPIURL: getEnv("TRIPS_API_URL", "http://localhost:8002"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),

		// RabbitMQ
		RabbitMQURL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQExchange: getEnv("RABBITMQ_EXCHANGE", "reservations.events"),
	}

	// Validate required fields
	if config.DBPassword == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
