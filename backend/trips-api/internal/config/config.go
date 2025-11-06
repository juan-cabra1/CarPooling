package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	ServerPort   string
	RabbitMQURL  string
	UsersAPIURL  string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	return &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:     getEnv("MONGO_DB", "carpooling_trips"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		ServerPort:  getEnv("SERVER_PORT", "8002"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
