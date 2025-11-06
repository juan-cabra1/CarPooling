package config

import (
	"fmt"
	"log"
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
// Retorna error si alguna variable crítica no está definida
func LoadConfig() (*Config, error) {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8002"),
		Mongo: MongoConfig{
			URI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DB:  getEnv("MONGO_DB", "carpooling_trips"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		},
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		UsersAPIURL: getEnv("USERS_API_URL", "http://localhost:8001"),
	}

	// Validar configuración crítica
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate verifica que la configuración sea válida
func (c *Config) Validate() error {
	if c.Mongo.URI == "" {
		return fmt.Errorf("MONGO_URI es requerido")
	}

	if c.Mongo.DB == "" {
		return fmt.Errorf("MONGO_DB es requerido")
	}

	if c.RabbitMQ.URL == "" {
		return fmt.Errorf("RABBITMQ_URL es requerido")
	}

	if c.JWTSecret == "" || c.JWTSecret == "your-super-secret-jwt-key-change-this-in-production" {
		log.Println("⚠️  ADVERTENCIA: Usando JWT_SECRET por defecto. Cambiar en producción!")
	}

	if c.UsersAPIURL == "" {
		return fmt.Errorf("USERS_API_URL es requerido")
	}

	return nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
