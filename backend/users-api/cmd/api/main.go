package main

import (
	"log"
	"users-api/internal/config"
	"users-api/internal/controller"
	"users-api/internal/dao"
	"users-api/internal/repository"
	"users-api/internal/routes"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// 2. Conectar a MySQL usando GORM
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos establecida")

	// 3. Auto-migrar los modelos (crear tablas si no existen)
	err = db.AutoMigrate(&dao.UserDAO{}, &dao.RatingDAO{})
	if err != nil {
		log.Fatalf("Error en auto-migración: %v", err)
	}

	log.Println("Auto-migración completada")

	// 4. Inicializar repositorios
	userRepo := repository.NewUserRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	// 5. Inicializar servicios
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(userRepo, emailService, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)
	ratingService := service.NewRatingService(ratingRepo, userRepo)

	// 6. Inicializar controladores
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	ratingController := controller.NewRatingController(ratingService)

	// 7. Crear router Gin
	router := gin.Default()

	// 8. Configurar rutas
	routes.SetupRoutes(router, authController, userController, ratingController, authService, userRepo)

	// 9. Iniciar servidor
	port := ":" + cfg.ServerPort
	log.Printf("Servidor iniciado en el puerto %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
