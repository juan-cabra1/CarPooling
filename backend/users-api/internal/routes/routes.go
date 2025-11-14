package routes

import (
	"users-api/internal/controller"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(
	router *gin.Engine,
	authController controller.AuthController,
	userController controller.UserController,
	ratingController controller.RatingController,
	authService service.AuthService,
	userRepo repository.UserRepository,
) {
	// Middleware globales
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())

	// ==================== HEALTH CHECK ====================
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data":    gin.H{"status": "ok"},
		})
	})

	// ==================== RUTAS PÚBLICAS (sin autenticación) ====================

	// Registro y Login
	router.POST("/users", authController.Register)
	router.POST("/login", authController.Login)

	// Verificación de email y recuperación de contraseña
	router.GET("/verify-email", authController.VerifyEmail)
	router.POST("/resend-verification", authController.ResendVerificationEmail)
	router.POST("/forgot-password", authController.RequestPasswordReset)
	router.POST("/reset-password", authController.ResetPassword)

	// ==================== RUTAS PROTEGIDAS (requieren JWT + Email verificado) ====================

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.Use(middleware.RequireVerifiedEmail(userRepo))
	{
		// Perfil de usuario
		protected.GET("/users/me", userController.GetMe)
		protected.GET("/users/:id", userController.GetUserByID)
		protected.PUT("/users/:id", userController.UpdateUser)
		protected.DELETE("/users/:id", userController.DeleteUser)

		// Calificaciones de usuario
		protected.GET("/users/:id/ratings", ratingController.GetUserRatings)

		// Cambio de contraseña
		protected.POST("/change-password", authController.ChangePassword)
	}

	// ==================== RUTAS INTERNAS (sin autenticación, para comunicación entre servicios) ====================

	internal := router.Group("/internal")
	{
		// Obtener usuario (llamado desde search-api y otros servicios)
		internal.GET("/users/:id", userController.GetUserByID)

		// Crear calificación (llamado desde trips-api)
		internal.POST("/ratings", ratingController.CreateRating)
	}
}
