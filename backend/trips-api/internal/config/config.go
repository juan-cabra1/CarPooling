package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	Mongo       MongoConfig
	RabbitMQ    RabbitMQConfig
	JWTSecret   string
	UsersAPIURL string
}

type MongoConfig struct {
	URI string
	DB  string
}

type RabbitMQConfig struct {
	URL string
}

// LoadConfig carga la configuración desde variables de entorno
// Usa fail-fast: panic si alguna variable crítica no está definida
func LoadConfig() (*Config, error) {
	// Intentar cargar .env desde la raíz del proyecto
	// En Docker, las variables vienen del docker-compose, así que esto falla silenciosamente
	_ = godotenv.Load()

	cfg := &Config{
		// Variables CRÍTICAS - Sin defaults, DEBEN existir (fail-fast)
		JWTSecret: mustGetEnv("JWT_SECRET"),
		Mongo: MongoConfig{
			URI: mustGetEnv("MONGO_URI_TRIPS"),
			DB:  getEnv("MONGO_DB_TRIPS", "carpooling_trips"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: mustGetEnv("RABBITMQ_URL"),
		},

		// Variables NO CRÍTICAS - Con defaults razonables
		ServerPort:  getEnv("SERVER_PORT", "8002"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
	}

	return cfg, nil
}

// getEnv obtiene variable con fallback (solo para variables NO críticas)
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
