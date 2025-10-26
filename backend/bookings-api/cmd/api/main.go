package main

import (
	"fmt"
	"log"

	"github.com/carpooling-ucc/bookings-api/internal/config"
	"github.com/carpooling-ucc/bookings-api/internal/controller"
	"github.com/carpooling-ucc/bookings-api/internal/dao"
	"github.com/carpooling-ucc/bookings-api/internal/middleware"
	"github.com/carpooling-ucc/bookings-api/internal/repository"
	"github.com/carpooling-ucc/bookings-api/internal/routes"
	"github.com/carpooling-ucc/bookings-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully")

	// 2. Connect to MySQL database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// Auto-migrate database schema
	if err := db.AutoMigrate(&dao.BookingDAO{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	// 3. Connect to RabbitMQ
	rabbitClient, err := service.NewRabbitMQClient(cfg.RabbitMQURL, cfg.RabbitMQExchange)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()
	log.Println("RabbitMQ connected successfully")

	// 4. Initialize HTTP clients for external services
	tripsClient := service.NewTripsClient(cfg.TripsAPIURL)
	usersClient := service.NewUsersClient(cfg.UsersAPIURL)
	log.Println("HTTP clients initialized")

	// 5. Initialize repository layer
	bookingRepo := repository.NewBookingRepository(db)
	log.Println("Repository layer initialized")

	// 6. Initialize service layer (with dependency injection)
	validationService := service.NewValidationService(tripsClient, usersClient, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, validationService, tripsClient, rabbitClient)
	log.Println("Service layer initialized")

	// 7. Initialize controller layer
	bookingController := controller.NewBookingController(bookingService)
	log.Println("Controller layer initialized")

	// 8. Setup Gin router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler())

	// 9. Setup routes with auth middleware
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	routes.SetupRoutes(router, bookingController, authMiddleware)
	log.Println("Routes configured successfully")

	// 10. Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("Starting bookings-api server on %s...", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
