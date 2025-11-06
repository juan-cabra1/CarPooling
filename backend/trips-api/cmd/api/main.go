package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trips-api/internal/clients"
	"trips-api/internal/config"
	"trips-api/internal/controllers"
	"trips-api/internal/middleware"
	"trips-api/internal/repository"
	"trips-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 📋 Cargar configuración desde las variables de entorno
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}
	log.Println("✅ Configuración cargada exitosamente")

	// 🗃️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// Context principal de la aplicación
	ctx := context.Background()

	// 🔌 Capa de datos: maneja operaciones con MongoDB
	tripsMongoRepo, err := repository.NewMongoTripsRepository(
		ctx,
		cfg.Mongo.URI,
		cfg.Mongo.DB,
		"trips",
	)
	if err != nil {
		log.Fatalf("Error inicializando repositorio de trips: %v", err)
	}
	log.Printf("✅ Conexión a MongoDB establecida [%s/%s]", cfg.Mongo.URI, cfg.Mongo.DB)

	// 📨 Inicializar RabbitMQ para publicar eventos de trips
	tripsQueue, err := clients.NewRabbitMQClient(
		cfg.RabbitMQ.URL,
		"trips.events", // Exchange name
		"topic",        // Exchange type
	)
	if err != nil {
		log.Fatalf("Error inicializando RabbitMQ: %v", err)
	}
	defer tripsQueue.Close()
	log.Println("✅ Conexión a RabbitMQ establecida")

	// 🌐 Cliente HTTP para comunicación con users-api
	usersAPIClient := clients.NewUsersAPIClient(cfg.UsersAPIURL, 10*time.Second)

	// 🔧 Capa de lógica de negocio: validaciones, transformaciones
	tripService := services.NewTripsService(
		tripsMongoRepo,
		tripsQueue,
		usersAPIClient,
	)

	// 🎮 Capa de controladores: maneja HTTP requests/responses
	tripController := controllers.NewTripsController(tripService)

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// Middleware: funciones que se ejecutan en cada request
	router.Use(middleware.CORSMiddleware())

	// 🏥 Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "trips-api",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 🚗 Rutas de Trips API
	// Grupo de rutas que requieren autenticación
	api := router.Group("/api/v1")
	{
		// POST /api/v1/trips - crear nuevo viaje (requiere auth)
		api.POST("/trips", middleware.AuthMiddleware(cfg.JWTSecret), tripController.CreateTrip)

		// GET /api/v1/trips/:id - obtener viaje por ID
		api.GET("/trips/:id", tripController.GetTripByID)

		// PUT /api/v1/trips/:id - actualizar viaje existente (requiere auth + ownership)
		api.PUT("/trips/:id", middleware.AuthMiddleware(cfg.JWTSecret), tripController.UpdateTrip)

		// DELETE /api/v1/trips/:id - eliminar viaje (requiere auth + ownership)
		api.DELETE("/trips/:id", middleware.AuthMiddleware(cfg.JWTSecret), tripController.DeleteTrip)

		// GET /api/v1/trips/user/:userId - obtener viajes de un usuario
		api.GET("/trips/user/:userId", tripController.GetTripsByUser)

		// POST /api/v1/trips/:id/reserve - endpoint de acción (delega a reservations-api)
		api.POST("/trips/:id/reserve", middleware.AuthMiddleware(cfg.JWTSecret), tripController.ReserveTrip)
	}

	// Configuración del server HTTP con timeouts
	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 🚀 Iniciar servidor en una goroutine
	go func() {
		log.Printf("🚀 Trips API listening on port %s", cfg.ServerPort)
		log.Printf("🏥 Health check: http://localhost:%s/health", cfg.ServerPort)
		log.Printf("🚗 Trips API: http://localhost:%s/api/v1/trips", cfg.ServerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando el servidor: %v", err)
		}
	}()

	// 🛑 Graceful shutdown
	// Esperar señal de interrupción (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️  Apagando servidor...")

	// Dar tiempo para que las conexiones terminen (máximo 5 segundos)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error en shutdown del servidor: %v", err)
	}

	log.Println("✅ Servidor detenido correctamente")
}
