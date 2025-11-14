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
	godotenv.Load()

	// Construir DATABASE_URL desde variables individuales si no existe
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		dbUser := getEnv("DB_USER", "root")
		dbPassword := getEnv("DB_PASSWORD", "")
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "3306")
		dbName := getEnv("DB_NAME", "carpooling_users")

		databaseURL = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return &Config{
		DatabaseURL:  databaseURL,
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:   getEnv("SERVER_PORT", "8001"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPFrom:     getEnv("SMTP_FROM", "marcelinho.nelson@gmail.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "nhrw ylah yhvw qraj"),
		AppURL:       getEnv("APP_URL", "http://localhost:3000"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
