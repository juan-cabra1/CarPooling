package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
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

	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "Prueba.9876"),
		DBName:       getEnv("DB_NAME", "carpooling_users"),
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
