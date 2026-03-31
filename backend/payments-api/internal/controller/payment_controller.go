package controller

import (
	"net/http"

	"payments-api/internal/domain"
	"payments-api/internal/middleware"
	"payments-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// PaymentController handles payment endpoints
type PaymentController struct {
	paymentService service.PaymentService
}

// NewPaymentController creates a new payment controller
func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// CreatePreference creates a payment preference for a booking
// @Summary Create payment preference
// @Description Creates a Mercado Pago payment preference for a booking
// @Tags payments
// @Accept json
// @Produce json
// @Param request body domain.CreatePreferenceRequest true "Preference request"
// @Success 201 {object} domain.CreatePreferenceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/payments/preference [post]
func (c *PaymentController) CreatePreference(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Get email from JWT claims — this is server-verified, not from the request body
	passengerEmail, _ := middleware.GetUserEmail(ctx)

	var req domain.CreatePreferenceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	resp, err := c.paymentService.CreatePreference(userID, passengerEmail, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create preference")

		status := http.StatusInternalServerError
		switch err {
		case domain.ErrPaymentAlreadyExists:
			status = http.StatusConflict
		case domain.ErrDriverNotFound, domain.ErrUserNotFound:
			status = http.StatusUnprocessableEntity
		case domain.ErrSellerNotActive:
			status = http.StatusUnprocessableEntity
		case domain.ErrInvalidRequest:
			status = http.StatusBadRequest
		}

		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetPayment retrieves a payment by UUID
// @Summary Get payment
// @Description Retrieves payment details by UUID
// @Tags payments
// @Produce json
// @Param id path string true "Payment UUID"
// @Success 200 {object} domain.PaymentResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/payments/{id} [get]
func (c *PaymentController) GetPayment(ctx *gin.Context) {
	paymentUUID := ctx.Param("id")

	resp, err := c.paymentService.GetPayment(paymentUUID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetPaymentByBooking retrieves a payment by booking ID
// @Summary Get payment by booking
// @Description Retrieves payment details by booking ID
// @Tags payments
// @Produce json
// @Param booking_id query string true "Booking UUID"
// @Success 200 {object} domain.PaymentResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/payments/by-booking [get]
func (c *PaymentController) GetPaymentByBooking(ctx *gin.Context) {
	bookingID := ctx.Query("booking_id")
	if bookingID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "booking_id es requerido",
		})
		return
	}

	resp, err := c.paymentService.GetPaymentByBooking(bookingID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// ReleasePayment releases held funds to driver
// @Summary Release payment
// @Description Releases held funds to driver's available balance
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment UUID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/payments/{id}/release [post]
func (c *PaymentController) ReleasePayment(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	paymentUUID := ctx.Param("id")

	err := c.paymentService.ReleasePayment(userID, paymentUUID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrPaymentNotFound:
			status = http.StatusNotFound
		case domain.ErrForbidden:
			status = http.StatusForbidden
		case domain.ErrCannotReleasePayment:
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Fondos liberados exitosamente",
	})
}

// Webhook handles Mercado Pago webhook notifications
// @Summary MP Webhook
// @Description Processes webhook notifications from Mercado Pago
// @Tags webhooks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/webhooks/mercadopago [post]
func (c *PaymentController) Webhook(ctx *gin.Context) {
	var payload domain.MPWebhookPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Warn().Err(err).Msg("Invalid webhook payload")
		// Still return 200 to avoid retries
		ctx.JSON(http.StatusOK, gin.H{"received": true})
		return
	}

	log.Info().
		Int64("id", payload.ID).
		Str("type", payload.Type).
		Str("action", payload.Action).
		Msg("Received MP webhook")

	// Process asynchronously to return quickly
	go func() {
		if err := c.paymentService.ProcessWebhook(&payload); err != nil {
			log.Error().Err(err).Msg("Failed to process webhook")
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{"received": true})
}
