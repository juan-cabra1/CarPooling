package routes

import (
	"payments-api/internal/controller"
	"payments-api/internal/middleware"
	"payments-api/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the payments-api
func SetupRoutes(
	router *gin.Engine,
	healthController *controller.HealthController,
	paymentController *controller.PaymentController,
	walletController *controller.WalletController,
	sellerController *controller.SellerController,
	authService service.AuthService,
) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorMiddleware())

	// Public routes (no authentication required)
	router.GET("/health", healthController.Health)
	router.GET("/ready", healthController.Ready)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Webhook routes (no auth - validated by signature)
	webhooks := v1.Group("/webhooks")
	{
		// Mercado Pago webhook - receives payment notifications
		webhooks.POST("/mercadopago", paymentController.Webhook)
	}

	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Payment routes
		payments := protected.Group("/payments")
		{
			// Create payment preference for a booking
			payments.POST("/preference", paymentController.CreatePreference)

			// Get payment by booking ID
			payments.GET("/by-booking", paymentController.GetPaymentByBooking)

			// Get payment by ID
			payments.GET("/:id", paymentController.GetPayment)

			// Release funds (passenger confirms trip completion)
			payments.POST("/:id/release", paymentController.ReleasePayment)
		}

		// Wallet routes
		wallet := protected.Group("/wallet")
		{
			// Get my wallet
			wallet.GET("", walletController.GetWallet)

			// Get wallet transactions
			wallet.GET("/transactions", walletController.GetTransactions)

			// Get wallet statistics
			wallet.GET("/stats", walletController.GetStats)

			// Update wallet settings (auto-withdraw)
			wallet.PATCH("/settings", walletController.UpdateSettings)
		}

		// Withdrawal routes
		withdrawals := protected.Group("/withdrawals")
		{
			// Request a withdrawal
			withdrawals.POST("", walletController.RequestWithdrawal)

			// Get my withdrawals
			withdrawals.GET("", walletController.GetWithdrawals)

			// Get withdrawal by ID
			withdrawals.GET("/:id", walletController.GetWithdrawal)
		}

		// Seller (driver MP account) routes
		seller := protected.Group("/seller")
		{
			// Link MP account manually with access token
			seller.POST("/link-manual", sellerController.LinkManualToken)

			// Get OAuth URL to link MP account
			seller.GET("/oauth/url", sellerController.GetOAuthURL)

			// OAuth callback (receives authorization code)
			seller.POST("/oauth/callback", sellerController.OAuthCallback)

			// Get seller account status
			seller.GET("/status", sellerController.GetStatus)

			// Disconnect MP account
			seller.DELETE("/disconnect", sellerController.Disconnect)
		}
	}
}
