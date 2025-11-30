package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	ServerPort   string
	SMTPHost     string
	SMTPPort     string
	SMTPFrom     string
	SMTPPassword string
	AppURL       string
}

func LoadConfig() (*Config, error) {
	// Intentar cargar .env desde la raíz del proyecto
	// En Docker, las variables vienen del docker-compose, así que esto falla silenciosamente
	_ = godotenv.Load()

	return &Config{
		// Variables CRÍTICAS - Sin defaults, DEBEN existir
		DatabaseURL:  buildDatabaseURL(),
		JWTSecret:    mustGetEnv("JWT_SECRET"),
		SMTPPassword: mustGetEnv("SMTP_PASSWORD"),
		AppURL:       mustGetEnv("APP_URL"),

		// Variables NO CRÍTICAS - Con defaults razonables
		ServerPort: getEnv("SERVER_PORT", "8001"),
		SMTPHost:   getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:   getEnv("SMTP_PORT", "587"),
		SMTPFrom:   getEnv("SMTP_FROM", "matiasjbocco@gmail.com"),
	}, nil
}

// buildDatabaseURL construye la URL desde variables individuales o usa DATABASE_URL directamente
func buildDatabaseURL() string {
	// Opción 1: DATABASE_URL completa (preferido en producción)
	if dbURL := os.Getenv("DATABASE_URL_USERS"); dbURL != "" {
		return dbURL
	}

	// Opción 2: Construir desde componentes (desarrollo)
	dbUser := mustGetEnv("DB_USER")
	dbPassword := mustGetEnv("DB_PASSWORD")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME_USERS", "carpooling_users")

	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
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
