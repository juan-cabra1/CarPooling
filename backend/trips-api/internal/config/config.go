package config

import (
	"fmt"
	"os"
)

type MongoConfig struct {
	URI string
	DB  string
}

type Config struct {
	Mongo MongoConfig
	Port  string
}

// LoadConfig carga la configuración desde variables de entorno y valida
func LoadConfig() (Config, error) {
	cfg := Config{
		Mongo: MongoConfig{
			URI: os.Getenv("MONGO_URI"),
			DB:  os.Getenv("MONGO_DB"),
		},
		Port: os.Getenv("PORT"),
	}

	if cfg.Mongo.URI == "" {
		return cfg, fmt.Errorf("MONGO_URI no configurada")
	}
	if cfg.Mongo.DB == "" {
		return cfg, fmt.Errorf("MONGO_DB no configurada")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}