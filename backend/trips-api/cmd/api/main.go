package main

import (
	"context"
	"log"
	"time"
	"trips-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Configurar zerolog para logging estructurado
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: gin.DefaultWriter, TimeFormat: time.RFC3339})

	// 1. Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	zlog.Info().Msg("Configuración cargada exitosamente")

	// 2. Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error conectando a MongoDB")
	}

	// Verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Error verificando conexión a MongoDB")
	}

	zlog.Info().Str("uri", cfg.MongoURI).Str("database", cfg.MongoDB).Msg("Conexión a MongoDB establecida")

	// Obtener referencia a la base de datos
	db := client.Database(cfg.MongoDB)
	_ = db // Usaremos la base de datos en fases posteriores

	// 3. Crear router Gin
	router := gin.Default()

	// 4. Configurar endpoint de health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "trips-api",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 5. Iniciar servidor
	port := ":" + cfg.ServerPort
	zlog.Info().Str("port", cfg.ServerPort).Msg("Servidor iniciado")
	if err := router.Run(port); err != nil {
		zlog.Fatal().Err(err).Msg("Error iniciando el servidor")
	}
}
