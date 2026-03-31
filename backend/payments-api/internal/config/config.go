package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the payments-api microservice
// It uses environment variables with sensible defaults for development
type Config struct {
	// ServerPort is the HTTP port the server will listen on
	// Default: "8005"
	ServerPort string

	// DatabaseURL is the MySQL connection string in DSN format
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// Example: "root:password@tcp(localhost:3306)/payments_db?charset=utf8mb4&parseTime=True&loc=Local"
	DatabaseURL string

	// JWTSecret is the secret key used for signing and validating JWT tokens
	// IMPORTANT: Must match the secret used by other microservices for authentication
	// Should be a strong, random string in production
	JWTSecret string

	// RabbitMQURL is the connection URL for RabbitMQ message broker
	// Format: amqp://user:password@host:port/
	// Example: "amqp://guest:guest@localhost:5672/"
	RabbitMQURL string

	// Environment specifies the deployment environment
	// Values: "development", "staging", "production"
	// Affects logging verbosity and error handling behavior
	Environment string

	// UsersAPIURL is the base URL for the users-api microservice
	// Format: http://host:port (no trailing slash)
	UsersAPIURL string

	// TripsAPIURL is the base URL for the trips-api microservice
	// Format: http://host:port (no trailing slash)
	TripsAPIURL string

	// BookingsAPIURL is the base URL for the bookings-api microservice
	// Format: http://host:port (no trailing slash)
	BookingsAPIURL string

	// ============================================
	// MERCADO PAGO CONFIGURATION
	// ============================================

	// MPAccessToken is the Mercado Pago access token for API calls
	// Get it from: https://www.mercadopago.com.ar/developers/panel/app
	MPAccessToken string

	// MPPublicKey is the Mercado Pago public key for frontend checkout
	MPPublicKey string

	// MPClientID is the Mercado Pago OAuth client ID
	// Used for seller account linking (OAuth flow)
	MPClientID string

	// MPClientSecret is the Mercado Pago OAuth client secret
	// Used for seller account linking (OAuth flow)
	MPClientSecret string

	// MPWebhookSecret is the secret used to validate webhook signatures from MP
	MPWebhookSecret string

	// ============================================
	// BUSINESS RULES CONFIGURATION
	// ============================================

	// PlatformFeePercentage is the commission percentage taken by the platform
	// Default: 15 (meaning 15% goes to platform, 85% to driver)
	PlatformFeePercentage int

	// AutoReleaseHours is the number of hours after trip completion
	// before funds are automatically released to the driver
	// Default: 2
	AutoReleaseHours int

	// EnableMarketplace enables marketplace mode with split payments
	// When true: Uses marketplace_fee in MP preferences (requires MP marketplace app)
	// When false: Regular payments without marketplace_fee (for testing)
	// Default: false
	EnableMarketplace bool

	// ============================================
	// ENCRYPTION CONFIGURATION
	// ============================================

	// EncryptionKey is the 32-byte key used for AES-256 encryption
	// Used to encrypt sensitive data like MP OAuth tokens
	// MUST be exactly 32 characters
	EncryptionKey string

	// GoogleMapsAPIKey is the Google Maps API key (server-side use only)
	GoogleMapsAPIKey string

	// FrontendWebURL is the base URL of the web frontend (e.g. https://carpooling.railway.app)
	// Used to build back_url for MercadoPago payment preferences.
	// If empty, mobile deep links (carpooling://) are used instead.
	FrontendWebURL string
}

// LoadConfig reads configuration from environment variables
// It first attempts to load from a .env file, then falls back to system environment variables
//
// Returns:
//   - *Config: Loaded configuration
//   - error: Error if required variables are missing or invalid
func LoadConfig() (*Config, error) {
	// Attempt to load .env file from current directory
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("Loaded configuration from .env file")
	}

	// Build configuration from environment variables with defaults
	cfg := &Config{
		// Server
		ServerPort:  getEnv("SERVER_PORT", "8005"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		RabbitMQURL: getEnv("RABBITMQ_URL", ""),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Microservice URLs
		UsersAPIURL:    getEnv("USERS_API_URL", "http://localhost:8001"),
		TripsAPIURL:    getEnv("TRIPS_API_URL", "http://localhost:8002"),
		BookingsAPIURL: getEnv("BOOKINGS_API_URL", "http://localhost:8003"),

		// Mercado Pago
		MPAccessToken:   getEnv("MP_ACCESS_TOKEN", ""),
		MPPublicKey:     getEnv("MP_PUBLIC_KEY", ""),
		MPClientID:      getEnv("MP_CLIENT_ID", ""),
		MPClientSecret:  getEnv("MP_CLIENT_SECRET", ""),
		MPWebhookSecret: getEnv("MP_WEBHOOK_SECRET", ""),

		// Business rules
		PlatformFeePercentage: getEnvAsInt("PLATFORM_FEE_PERCENTAGE", 15),
		AutoReleaseHours:      getEnvAsInt("AUTO_RELEASE_HOURS", 2),
		EnableMarketplace:     getEnvAsBool("ENABLE_MARKETPLACE", false),

		// Encryption
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),

		// Maps & Frontend
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		FrontendWebURL:   getEnv("FRONTEND_WEB_URL", ""),
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration fields are set
// Returns an error if any required field is missing or invalid
func (c *Config) Validate() error {
	// Check DatabaseURL is set
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Check RabbitMQURL is set
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}

	// Check JWTSecret is set (critical for security)
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	// Check Mercado Pago credentials
	if c.MPAccessToken == "" {
		return fmt.Errorf("MP_ACCESS_TOKEN is required")
	}

	// Validate platform fee is reasonable (0-100%)
	if c.PlatformFeePercentage < 0 || c.PlatformFeePercentage > 100 {
		return fmt.Errorf("PLATFORM_FEE_PERCENTAGE must be between 0 and 100")
	}

	// Validate auto-release hours
	if c.AutoReleaseHours < 0 {
		return fmt.Errorf("AUTO_RELEASE_HOURS must be non-negative")
	}

	// Warn if using development environment in production
	if c.Environment == "production" {
		if c.JWTSecret == "dev-secret" {
			log.Println("WARNING: Using default JWT_SECRET in production is insecure!")
		}
		if c.EncryptionKey == "" {
			log.Println("WARNING: ENCRYPTION_KEY not set - OAuth tokens will not be encrypted!")
		}
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetDriverAmountPercentage returns the percentage that goes to the driver
// This is calculated as 100 - platform fee
func (c *Config) GetDriverAmountPercentage() int {
	return 100 - c.PlatformFeePercentage
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as an integer with a fallback default
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: %s is not a valid integer, using default %d", key, defaultValue)
		return defaultValue
	}
	return intValue
}

// getEnvAsBool retrieves an environment variable as a boolean with a fallback default
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: %s is not a valid boolean, using default %v", key, defaultValue)
		return defaultValue
	}
	return boolValue
}
