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
	verificationController controller.VerificationController,
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
	router.POST("/users", middleware.RateLimitMiddleware(20), authController.Register)
	router.POST("/login", middleware.RateLimitMiddleware(10), authController.Login)

	// Verificación de email y recuperación de contraseña
	router.GET("/verify-email", authController.VerifyEmail)
	router.POST("/resend-verification", middleware.RateLimitMiddleware(5), authController.ResendVerificationEmail)
	router.POST("/forgot-password", middleware.RateLimitMiddleware(5), authController.RequestPasswordReset)
	router.POST("/reset-password", middleware.RateLimitMiddleware(10), authController.ResetPassword)

	// ==================== RUTAS PROTEGIDAS (requieren JWT + Email verificado) ====================

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.Use(middleware.RequireVerifiedEmail(userRepo))
	protected.Use(middleware.RateLimitMiddleware(120))
	{
		// Perfil de usuario
		protected.GET("/users/me", userController.GetMe)
		protected.GET("/users/:id", userController.GetUserByID)
		protected.PUT("/users/:id", userController.UpdateUser)
		protected.DELETE("/users/:id", userController.DeleteUser)

		// Calificaciones de usuario
		protected.GET("/users/:id/ratings", ratingController.GetUserRatings)

		// Verificación de identidad
		protected.POST("/users/:id/verification", verificationController.SubmitVerification)
		protected.GET("/users/:id/verification", verificationController.GetVerificationStatus)

		// Cambio de contraseña
		protected.POST("/change-password", authController.ChangePassword)
	}

	// ==================== RUTAS ADMIN (requieren JWT + rol admin) ====================

	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.RequireAdminRole())
	{
		// Gestión de usuarios (solo admin)
		admin.GET("/users", userController.GetAllUsers)
		admin.POST("/users/:id/force-reauth", userController.ForceReauthentication)
		admin.POST("/users/:id/block", userController.BlockUser)
		admin.POST("/users/:id/unblock", userController.UnblockUser)

		// Verificaciones de identidad (solo admin)
		admin.GET("/verifications", verificationController.GetPendingVerifications)
		admin.POST("/verifications/:id/review", verificationController.ReviewVerification)
	}

	// ==================== RUTAS INTERNAS (requieren X-Internal-Token, para comunicación entre servicios) ====================

	internal := router.Group("/internal")
	internal.Use(middleware.InternalAuthMiddleware())
	{
		// Obtener usuario (llamado desde search-api y otros servicios)
		internal.GET("/users/:id", userController.GetUserByID)

		// Crear calificación (llamado desde trips-api)
		internal.POST("/ratings", ratingController.CreateRating)
	}
}
