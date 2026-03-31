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
	"trips-api/internal/controller"
	"trips-api/internal/database"
	"trips-api/internal/messaging"
	"trips-api/internal/middleware"
	"trips-api/internal/repository"
	"trips-api/internal/routes"
	"trips-api/internal/service"
	ws "trips-api/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	// 🔧 Configurar zerolog para logging estructurado en JSON
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 📋 Cargar configuración desde las variables de entorno
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}
	log.Println("✅ Configuración cargada exitosamente")

	// 🗃️ Inicializar capas de la aplicación (Dependency Injection)
	// Patrón: Repository -> Service -> Controller
	// Cada capa tiene una responsabilidad específica

	// 🔌 Conectar a MongoDB
	db, err := database.ConnectMongoDB(cfg.Mongo.URI, cfg.Mongo.DB)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}

	// 📋 Crear colecciones e índices
	if err := database.CreateIndexes(db); err != nil {
		log.Fatalf("Error creando índices: %v", err)
	}

	// 🔌 Capa de datos: maneja operaciones con MongoDB
	tripsRepo := repository.NewTripRepository(db)
	eventsRepo := repository.NewEventRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	log.Println("✅ Repositories initialized")

	// 🌐 Capa de clientes HTTP externos
	usersClient := clients.NewUsersClient(cfg.UsersAPIURL)
	log.Println("✅ HTTP clients initialized")

	// 📨 Conectar a RabbitMQ
	publisher, err := messaging.NewPublisher(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer publisher.Close()
	log.Println("✅ RabbitMQ publisher initialized")

	// 🔌 Iniciar WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()
	log.Println("✅ WebSocket hub started")

	// 📦 Capa de servicios: lógica de negocio
	idempotencyService := service.NewIdempotencyService(eventsRepo)
	tripService := service.NewTripService(tripsRepo, idempotencyService, usersClient, publisher)
	chatService := service.NewChatService(messageRepo, tripsRepo, publisher, wsHub)
	trackingService := service.NewTrackingService(tripsRepo)
	log.Println("✅ Services initialized")

	// 📥 Inicializar RabbitMQ consumer
	consumer, err := messaging.NewReservationConsumer(
		cfg.RabbitMQ.URL,
		tripService,
		idempotencyService,
		publisher,
	)
	if err != nil {
		log.Fatalf("Error inicializando consumer: %v", err)
	}
	defer consumer.Close()
	log.Println("✅ RabbitMQ consumer initialized")

	// Crear contexto para consumer (con cancelación)
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	// 🚀 Iniciar consumer en goroutine
	go func() {
		log.Println("🔄 Starting reservation events consumer...")
		if err := consumer.Start(consumerCtx); err != nil {
			log.Printf("⚠️  Consumer stopped: %v", err)
		}
	}()

	// 🎮 Capa de controladores: HTTP handlers
	authService := service.NewAuthService(cfg.JWTSecret)
	tripController := controller.NewTripController(tripService)
	chatController := controller.NewChatController(chatService)
	trackingController := controller.NewTrackingController(trackingService)
	wsController := controller.NewWSController(wsHub, cfg.JWTSecret)
	websocketController := controller.NewWebSocketController(wsHub, tripService)
	routeController := controller.NewRouteController(cfg.GoogleMapsAPIKey)
	log.Println("✅ Controllers initialized")

	// 🌐 Configurar router HTTP con Gin
	router := gin.Default()

	// 🌍 CORS middleware — debe aplicarse antes que cualquier ruta
	router.Use(middleware.CORSMiddleware())

	// 🔐 Crear JWT middleware
	jwtMiddleware := middleware.AuthMiddleware(authService)

	// 🚦 Configurar rutas de la aplicación
	routes.SetupRoutes(router, tripController, chatController, trackingController, wsController, websocketController, routeController, jwtMiddleware)
	log.Println("✅ Routes configured")

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
		log.Printf("🚗 Trips API: http://localhost:%s/trips", cfg.ServerPort)

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

	// Cerrar consumer de RabbitMQ
	log.Println("⚠️  Stopping consumer...")
	consumerCancel() // Signal consumer to stop
	time.Sleep(1 * time.Second) // Give consumer time to finish in-flight messages

	if err := consumer.Close(); err != nil {
		log.Printf("⚠️  Error closing consumer: %v", err)
	}

	// Cerrar publisher de RabbitMQ
	if err := publisher.Close(); err != nil {
		log.Printf("⚠️  Error cerrando RabbitMQ: %v", err)
	}

	log.Println("✅ Servidor detenido correctamente")
}
