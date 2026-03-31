package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	Mongo            MongoConfig
	RabbitMQ         RabbitMQConfig
	JWTSecret        string
	UsersAPIURL      string
	GoogleMapsAPIKey string // Para endpoint /route server-side
	FrontendWebURL   string // URL base del frontend web (ej: https://carpooling.railway.app)
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
		ServerPort:       getEnvFallback("PORT", "SERVER_PORT", "8002"),
		UsersAPIURL:      getEnv("USERS_API_URL", "http://localhost:8001"),
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		FrontendWebURL:   getEnv("FRONTEND_WEB_URL", ""),
	}

	return cfg, nil
}

// getEnvFallback intenta primary primero, luego secondary, luego el default
// Usado para compatibilidad con Railway (inyecta PORT) y desarrollo local (SERVER_PORT)
func getEnvFallback(primary, secondary, defaultValue string) string {
	if v := os.Getenv(primary); v != "" {
		return v
	}
	if v := os.Getenv(secondary); v != "" {
		return v
	}
	return defaultValue
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
