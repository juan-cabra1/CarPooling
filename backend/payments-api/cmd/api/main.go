package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payments-api/internal/clients"
	"payments-api/internal/config"
	"payments-api/internal/controller"
	"payments-api/internal/database"
	"payments-api/internal/repository"
	"payments-api/internal/routes"
	"payments-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// ========================================================================
	// LOGGING SETUP
	// ========================================================================
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	log.Info().Msg("========================================")
	log.Info().Msg("  PAYMENTS-API Starting...")
	log.Info().Msg("========================================")

	// ========================================================================
	// CONFIGURATION
	// ========================================================================
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().
		Str("port", cfg.ServerPort).
		Str("environment", cfg.Environment).
		Int("platform_fee_percentage", cfg.PlatformFeePercentage).
		Int("auto_release_hours", cfg.AutoReleaseHours).
		Msg("Configuration loaded")

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// ========================================================================
	// DATABASE
	// ========================================================================
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Run auto-migration
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run auto-migration")
	}

	// ========================================================================
	// EXTERNAL CLIENTS
	// ========================================================================
	mpClient := clients.NewMercadoPagoClient(cfg)
	log.Info().Msg("Mercado Pago client initialized")

	usersClient := clients.NewUsersClient(cfg.UsersAPIURL)
	log.Info().Str("users_api_url", cfg.UsersAPIURL).Msg("Users API client initialized")

	// ========================================================================
	// REPOSITORIES
	// ========================================================================
	paymentRepo := repository.NewPaymentRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	sellerRepo := repository.NewSellerRepository(db)
	withdrawalRepo := repository.NewWithdrawalRepository(db)
	eventRepo := repository.NewEventRepository(db)

	log.Info().Msg("Repositories initialized")

	// ========================================================================
	// SERVICES
	// ========================================================================
	authService := service.NewAuthService(cfg.JWTSecret)
	paymentService := service.NewPaymentService(cfg, db, paymentRepo, walletRepo, sellerRepo, eventRepo, mpClient, usersClient)
	walletService := service.NewWalletService(db, walletRepo)
	sellerService := service.NewSellerService(cfg, sellerRepo, mpClient)
	withdrawalService := service.NewWithdrawalService(cfg, withdrawalRepo, walletRepo, sellerRepo, mpClient)

	log.Info().Msg("Services initialized")

	// ========================================================================
	// CONTROLLERS
	// ========================================================================
	healthController := controller.NewHealthController(db)
	paymentController := controller.NewPaymentController(paymentService)
	walletController := controller.NewWalletController(walletService, withdrawalService)
	sellerController := controller.NewSellerController(sellerService)

	log.Info().Msg("Controllers initialized")

	// ========================================================================
	// ROUTER
	// ========================================================================
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add request logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Msg("Request")
	})

	// Setup routes
	routes.SetupRoutes(
		router,
		healthController,
		paymentController,
		walletController,
		sellerController,
		authService,
	)

	log.Info().Msg("Routes configured")

	// ========================================================================
	// BACKGROUND JOBS
	// ========================================================================
	// Start token refresh job
	go startTokenRefreshJob(sellerService)

	// Start auto-release job
	go startAutoReleaseJob(paymentService)

	// Start pending withdrawals job
	go startWithdrawalProcessorJob(withdrawalService)

	log.Info().Msg("Background jobs started")

	// ========================================================================
	// HTTP SERVER
	// ========================================================================
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().
			Str("port", cfg.ServerPort).
			Msg("HTTP server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	log.Info().Msg("========================================")
	log.Info().Msgf("  PAYMENTS-API Running on port %s", cfg.ServerPort)
	log.Info().Msg("========================================")

	// ========================================================================
	// GRACEFUL SHUTDOWN
	// ========================================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	// Close database connection
	if err := database.CloseDB(db); err != nil {
		log.Error().Err(err).Msg("Error closing database")
	}

	log.Info().Msg("Server exited")
}

// startTokenRefreshJob refreshes MP tokens every 5 minutes
func startTokenRefreshJob(sellerService service.SellerService) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Debug().Msg("Running token refresh job")
		if err := sellerService.RefreshExpiredTokens(); err != nil {
			log.Error().Err(err).Msg("Token refresh job failed")
		}
	}
}

// startAutoReleaseJob releases funds automatically after configured hours
func startAutoReleaseJob(paymentService service.PaymentService) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Debug().Msg("Running auto-release job")
		if err := paymentService.ProcessAutoRelease(); err != nil {
			log.Error().Err(err).Msg("Auto-release job failed")
		}
	}
}

// startWithdrawalProcessorJob processes pending withdrawals
func startWithdrawalProcessorJob(withdrawalService service.WithdrawalService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Debug().Msg("Running withdrawal processor job")
		if err := withdrawalService.ProcessPendingWithdrawals(); err != nil {
			log.Error().Err(err).Msg("Withdrawal processor job failed")
		}
	}
}
