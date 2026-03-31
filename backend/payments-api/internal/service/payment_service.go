package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"payments-api/internal/clients"
	"payments-api/internal/config"
	"payments-api/internal/dao"
	"payments-api/internal/domain"
	"payments-api/internal/repository"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// PaymentService handles payment business logic
type PaymentService interface {
	// CreatePreference creates a Mercado Pago preference for a booking.
	// passengerEmail must come from the authenticated user's JWT claims — never from the request body.
	CreatePreference(userID, passengerEmail string, req *domain.CreatePreferenceRequest) (*domain.CreatePreferenceResponse, error)
	GetPayment(paymentUUID string) (*domain.PaymentResponse, error)
	GetPaymentByBooking(bookingID string) (*domain.PaymentResponse, error)
	ProcessWebhook(payload *domain.MPWebhookPayload) error
	ReleasePayment(userID, paymentUUID string) error
	RefundPayment(paymentUUID string, percentage float64, reason string) error
	ProcessAutoRelease() error
}

type paymentService struct {
	cfg         *config.Config
	db          *gorm.DB
	paymentRepo repository.PaymentRepository
	walletRepo  repository.WalletRepository
	sellerRepo  repository.SellerRepository
	eventRepo   repository.EventRepository
	mpClient    clients.MercadoPagoClient
	usersClient clients.UsersClient
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	cfg *config.Config,
	db *gorm.DB,
	paymentRepo repository.PaymentRepository,
	walletRepo repository.WalletRepository,
	sellerRepo repository.SellerRepository,
	eventRepo repository.EventRepository,
	mpClient clients.MercadoPagoClient,
	usersClient clients.UsersClient,
) PaymentService {
	return &paymentService{
		cfg:         cfg,
		db:          db,
		paymentRepo: paymentRepo,
		walletRepo:  walletRepo,
		sellerRepo:  sellerRepo,
		eventRepo:   eventRepo,
		mpClient:    mpClient,
		usersClient: usersClient,
	}
}

// CreatePreference creates a payment preference for a booking.
// passengerEmail must be obtained from the JWT claims by the controller — not from the request body.
func (s *paymentService) CreatePreference(userID, passengerEmail string, req *domain.CreatePreferenceRequest) (*domain.CreatePreferenceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate that the driver exists in users-api before processing any payment.
	// DriverID comes from the frontend; we verify it's a real registered user.
	driverIDInt, err := strconv.ParseInt(req.DriverID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("driver_id", req.DriverID).Msg("Invalid driver_id format")
		return nil, domain.ErrInvalidRequest
	}

	if _, err := s.usersClient.GetUser(ctx, driverIDInt); err != nil {
		if err == domain.ErrUserNotFound {
			log.Warn().Str("driver_id", req.DriverID).Msg("Driver not found in users-api")
			return nil, domain.ErrDriverNotFound
		}
		// users-api is unavailable — fail closed to protect payment integrity
		log.Error().Err(err).Str("driver_id", req.DriverID).Msg("Failed to validate driver against users-api")
		return nil, domain.ErrInternalServer
	}

	// Check if driver has linked MercadoPago account
	seller, err := s.sellerRepo.FindActiveByUserID(req.DriverID)
	if err != nil {
		log.Error().Err(err).Str("driver_id", req.DriverID).Msg("Driver doesn't have active MercadoPago account")
		return nil, domain.ErrSellerNotActive
	}

	// Check if payment already exists for this booking
	existing, err := s.paymentRepo.FindByBookingID(req.BookingID)
	if err == nil && existing != nil {
		// Return existing preference if still pending
		if existing.IsPending() {
			return &domain.CreatePreferenceResponse{
				PaymentUUID:      existing.PaymentUUID,
				PreferenceID:     existing.MPPreferenceID,
				Amount:           existing.Amount,
				PlatformFee:      existing.PlatformFee,
				DriverAmount:     existing.DriverAmount,
				ExpiresAt:        *existing.ExpiresAt,
			}, nil
		}
		return nil, domain.ErrPaymentAlreadyExists
	}

	// Calculate amounts
	totalAmount := req.PricePerSeat * int64(req.SeatsCount)
	platformFee := (totalAmount * int64(s.cfg.PlatformFeePercentage)) / 100
	driverAmount := totalAmount - platformFee

	// Create payment record first
	expiresAt := time.Now().Add(30 * time.Minute)
	payment := &dao.Payment{
		BookingID:    req.BookingID,
		TripID:       req.TripID,
		PassengerID:  userID,
		DriverID:     req.DriverID,
		Amount:       totalAmount,
		Currency:     "ARS",
		PlatformFee:  platformFee,
		DriverAmount: driverAmount,
		SeatsCount:   req.SeatsCount,
		PricePerSeat: req.PricePerSeat,
		Status:       dao.PaymentStatusPending,
		FundStatus:   dao.FundStatusHeld,
		ExpiresAt:    &expiresAt,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		log.Error().Err(err).Msg("Failed to create payment record")
		return nil, domain.ErrInternalServer
	}

	// Create MP preference - using minimal fields like payment-service for compatibility
	mpReq := &domain.MPPreferenceRequest{
		Items: []domain.MPPreferenceItem{
			{
				Title:      fmt.Sprintf("Viaje %s → %s", req.Origin, req.Destination),
				Quantity:   1,
				CurrencyID: "ARS",
				UnitPrice:  math.Round(float64(totalAmount)/100*100) / 100, // Convert centavos to pesos and round to 2 decimal places
			},
		},
		ExternalReference: payment.PaymentUUID,
	}

	// Add marketplace configuration only if enabled
	// Marketplace mode requires MP app configured as marketplace
	if s.cfg.EnableMarketplace {
		mpReq.Marketplace = "CARPOOLING"
		mpReq.MarketplaceFee = math.Round(float64(platformFee)/100*100) / 100 // Platform fee in pesos, rounded
		log.Info().
			Str("payment_uuid", payment.PaymentUUID).
			Float64("marketplace_fee", mpReq.MarketplaceFee).
			Msg("Marketplace mode enabled - using marketplace_fee")
	} else {
		log.Info().
			Str("payment_uuid", payment.PaymentUUID).
			Msg("Marketplace mode disabled - using standard payment")
	}

	// Use passenger email from JWT claims (server-verified), not from the request body.
	// Fall back to req.PayerEmail only if JWT email is somehow empty (shouldn't happen).
	payerEmail := passengerEmail
	if payerEmail == "" {
		payerEmail = req.PayerEmail
	}
	if payerEmail != "" {
		mpReq.Payer = &domain.MPPayer{
			Email: payerEmail,
			Name:  req.PayerName,
		}
	}

	// Back URLs: el cliente decide el esquema según su plataforma.
	// Mobile envía "carpooling:/" → deep link que abre la app.
	// Web envía "https://carpooling.railway.app" → URL del navegador.
	// Fallback al default del servidor si el cliente no especifica.
	baseURL := req.ReturnBaseURL
	if baseURL == "" {
		baseURL = s.cfg.FrontendWebURL
	}
	if baseURL == "" {
		baseURL = "carpooling:/"
	}
	mpReq.BackURLs = &domain.MPBackURLs{
		Success: fmt.Sprintf("%s/payment/success?payment_id=%s", baseURL, payment.PaymentUUID),
		Pending: fmt.Sprintf("%s/payment/pending?payment_id=%s", baseURL, payment.PaymentUUID),
		Failure: fmt.Sprintf("%s/payment/failure?payment_id=%s", baseURL, payment.PaymentUUID),
	}
	mpReq.AutoReturn = "approved"

	// Choose which token to use based on marketplace mode
	// Marketplace mode: Use seller's token (split payment)
	// Standard mode: Use platform token (regular payment)
	var mpResp *domain.MPPreferenceResponse
	var mpErr error

	if s.cfg.EnableMarketplace {
		// Use seller's access token for marketplace payments
		log.Debug().Str("payment_uuid", payment.PaymentUUID).Msg("Using seller token for marketplace payment")
		mpResp, mpErr = s.mpClient.CreatePreferenceWithToken(mpReq, seller.MPAccessToken)
	} else {
		// Use platform token for standard payments
		log.Debug().Str("payment_uuid", payment.PaymentUUID).Msg("Using platform token for standard payment")
		mpResp, mpErr = s.mpClient.CreatePreference(mpReq)
	}

	if mpErr != nil {
		log.Error().Err(mpErr).Str("payment_uuid", payment.PaymentUUID).Msg("Failed to create MP preference")
		// Update payment status to error
		payment.Status = dao.PaymentStatusRejected
		s.paymentRepo.Update(payment)
		return nil, fmt.Errorf("%w: %v", domain.ErrMPPreferenceCreation, mpErr)
	}

	// Update payment with MP preference ID
	payment.MPPreferenceID = mpResp.ID
	if err := s.paymentRepo.Update(payment); err != nil {
		log.Error().Err(err).Msg("Failed to update payment with preference ID")
	}

	log.Info().
		Str("payment_uuid", payment.PaymentUUID).
		Str("preference_id", mpResp.ID).
		Int64("amount", totalAmount).
		Msg("Payment preference created")

	return &domain.CreatePreferenceResponse{
		PaymentUUID:      payment.PaymentUUID,
		PreferenceID:     mpResp.ID,
		InitPoint:        mpResp.InitPoint,
		SandboxInitPoint: mpResp.SandboxInitPoint,
		Amount:           totalAmount,
		PlatformFee:      platformFee,
		DriverAmount:     driverAmount,
		ExpiresAt:        expiresAt,
	}, nil
}

// GetPayment retrieves payment details
func (s *paymentService) GetPayment(paymentUUID string) (*domain.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByUUID(paymentUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return s.toPaymentResponse(payment), nil
}

// GetPaymentByBooking retrieves payment by booking ID
func (s *paymentService) GetPaymentByBooking(bookingID string) (*domain.PaymentResponse, error) {
	payment, err := s.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return s.toPaymentResponse(payment), nil
}

// ProcessWebhook processes a webhook notification from MP
func (s *paymentService) ProcessWebhook(payload *domain.MPWebhookPayload) error {
	if payload.Type != "payment" {
		log.Debug().Str("type", payload.Type).Msg("Ignoring non-payment webhook")
		return nil
	}

	paymentID := payload.Data.ID
	if paymentID == "" {
		return domain.ErrMPWebhookInvalid
	}

	// Get payment details from MP
	mpPayment, err := s.mpClient.GetPayment(paymentID)
	if err != nil {
		log.Error().Err(err).Str("mp_payment_id", paymentID).Msg("Failed to get payment from MP")
		return err
	}

	// Find our payment record
	payment, err := s.paymentRepo.FindByUUID(mpPayment.ExternalReference)
	if err != nil {
		log.Error().Err(err).Str("external_reference", mpPayment.ExternalReference).Msg("Payment not found")
		return domain.ErrPaymentNotFound
	}

	// Update payment with MP data
	payment.MPPaymentID = paymentID
	payment.MPStatus = mpPayment.Status
	payment.MPStatusDetail = mpPayment.StatusDetail
	payment.MPMerchantOrderID = strconv.FormatInt(mpPayment.MerchantOrderID, 10)

	// Update status based on MP status
	switch mpPayment.Status {
	case "approved":
		payment.Status = dao.PaymentStatusApproved
		now := time.Now()
		payment.ApprovedAt = &now

		// Credit driver's wallet with held funds
		if err := s.creditDriverWallet(payment); err != nil {
			log.Error().Err(err).Msg("Failed to credit driver wallet")
		}

	case "pending", "in_process":
		payment.Status = dao.PaymentStatusInProcess

	case "rejected":
		payment.Status = dao.PaymentStatusRejected

	case "cancelled":
		payment.Status = dao.PaymentStatusCancelled

	case "refunded":
		payment.Status = dao.PaymentStatusRefunded
		payment.FundStatus = dao.FundStatusRefunded
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		log.Error().Err(err).Msg("Failed to update payment")
		return err
	}

	log.Info().
		Str("payment_uuid", payment.PaymentUUID).
		Str("mp_status", mpPayment.Status).
		Msg("Payment webhook processed")

	return nil
}

// ReleasePayment releases held funds to driver's available balance
func (s *paymentService) ReleasePayment(userID, paymentUUID string) error {
	payment, err := s.paymentRepo.FindByUUID(paymentUUID)
	if err != nil {
		return domain.ErrPaymentNotFound
	}

	// Verify the requester is the passenger
	if payment.PassengerID != userID {
		return domain.ErrForbidden
	}

	if !payment.CanBeReleased() {
		return domain.ErrCannotReleasePayment
	}

	// Release funds in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get driver's wallet
		wallet, err := s.walletRepo.FindByUserID(payment.DriverID)
		if err != nil {
			return err
		}

		// Move from held to available
		wallet.HeldBalance -= payment.DriverAmount
		wallet.AvailableBalance += payment.DriverAmount
		wallet.TotalTrips++
		now := time.Now()
		wallet.LastDepositAt = &now

		if err := s.walletRepo.UpdateWithTx(tx, wallet); err != nil {
			return err
		}

		// Create transaction record
		transaction := &dao.WalletTransaction{
			WalletID:      wallet.ID,
			UserID:        payment.DriverID,
			Type:          dao.TransactionTypeRelease,
			Amount:        payment.DriverAmount,
			Currency:      "ARS",
			BalanceBefore: wallet.AvailableBalance - payment.DriverAmount,
			BalanceAfter:  wallet.AvailableBalance,
			PaymentID:     &payment.PaymentUUID,
			BookingID:     &payment.BookingID,
			TripID:        &payment.TripID,
			Description:   "Fondos liberados - viaje completado",
		}

		if err := s.walletRepo.CreateTransactionWithTx(tx, transaction); err != nil {
			return err
		}

		// Update payment status
		payment.FundStatus = dao.FundStatusReleased
		payment.ReleasedAt = &now

		return tx.Save(payment).Error
	})
}

// RefundPayment refunds a payment (full or partial)
func (s *paymentService) RefundPayment(paymentUUID string, percentage float64, reason string) error {
	payment, err := s.paymentRepo.FindByUUID(paymentUUID)
	if err != nil {
		return domain.ErrPaymentNotFound
	}

	if !payment.CanBeRefunded() {
		return domain.ErrCannotRefundPayment
	}

	// Default to full refund
	if percentage <= 0 || percentage > 100 {
		percentage = 100
	}

	refundAmount := float64(payment.Amount) * (percentage / 100)

	// Call MP refund API
	_, err = s.mpClient.CreateRefund(0, &refundAmount) // TODO: Get MP payment ID as int64
	if err != nil {
		log.Error().Err(err).Msg("Failed to create refund in MP")
		return domain.ErrMPRefundFailed
	}

	// Update payment
	payment.RefundedAmount = int64(refundAmount)
	payment.RefundReason = reason
	now := time.Now()
	payment.RefundedAt = &now

	if percentage == 100 {
		payment.Status = dao.PaymentStatusRefunded
		payment.FundStatus = dao.FundStatusRefunded
	} else {
		payment.Status = dao.PaymentStatusPartialRefund
	}

	return s.paymentRepo.Update(payment)
}

// creditDriverWallet adds funds to driver's held balance
func (s *paymentService) creditDriverWallet(payment *dao.Payment) error {
	wallet, err := s.walletRepo.GetOrCreate(payment.DriverID)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		wallet.HeldBalance += payment.DriverAmount
		wallet.TotalEarned += payment.DriverAmount
		now := time.Now()
		wallet.LastDepositAt = &now

		if err := s.walletRepo.UpdateWithTx(tx, wallet); err != nil {
			return err
		}

		// Create hold transaction
		transaction := &dao.WalletTransaction{
			WalletID:      wallet.ID,
			UserID:        payment.DriverID,
			Type:          dao.TransactionTypeHold,
			Amount:        payment.DriverAmount,
			Currency:      "ARS",
			BalanceBefore: wallet.HeldBalance - payment.DriverAmount,
			BalanceAfter:  wallet.HeldBalance,
			PaymentID:     &payment.PaymentUUID,
			BookingID:     &payment.BookingID,
			TripID:        &payment.TripID,
			Description:   "Pago recibido - fondos retenidos",
		}

		return s.walletRepo.CreateTransactionWithTx(tx, transaction)
	})
}

// toPaymentResponse converts a payment DAO to response
func (s *paymentService) toPaymentResponse(payment *dao.Payment) *domain.PaymentResponse {
	return &domain.PaymentResponse{
		PaymentUUID:  payment.PaymentUUID,
		BookingID:    payment.BookingID,
		TripID:       payment.TripID,
		PassengerID:  payment.PassengerID,
		DriverID:     payment.DriverID,
		Amount:       payment.Amount,
		Currency:     payment.Currency,
		PlatformFee:  payment.PlatformFee,
		DriverAmount: payment.DriverAmount,
		SeatsCount:   payment.SeatsCount,
		Status:       string(payment.Status),
		FundStatus:   string(payment.FundStatus),
		MPPaymentID:  payment.MPPaymentID,
		ApprovedAt:   payment.ApprovedAt,
		ReleasedAt:   payment.ReleasedAt,
		CreatedAt:    payment.CreatedAt,
	}
}

// ProcessAutoRelease releases funds for payments that are approved and past the auto-release window
func (s *paymentService) ProcessAutoRelease() error {
	// Find payments ready for auto-release
	payments, err := s.paymentRepo.FindPendingReleases(s.cfg.AutoReleaseHours)
	if err != nil {
		return err
	}

	for _, payment := range payments {
		if err := s.autoReleasePayment(&payment); err != nil {
			log.Error().
				Err(err).
				Str("payment_uuid", payment.PaymentUUID).
				Msg("Failed to auto-release payment")
		}
	}

	return nil
}

// autoReleasePayment releases a single payment's funds
func (s *paymentService) autoReleasePayment(payment *dao.Payment) error {
	if !payment.CanBeReleased() {
		return nil
	}

	// Get driver's wallet
	wallet, err := s.walletRepo.FindByUserID(payment.DriverID)
	if err != nil {
		return err
	}

	// Move from held to available
	wallet.HeldBalance -= payment.DriverAmount
	wallet.AvailableBalance += payment.DriverAmount
	wallet.TotalTrips++
	now := time.Now()
	wallet.LastDepositAt = &now

	if err := s.walletRepo.Update(wallet); err != nil {
		return err
	}

	// Create transaction record
	transaction := &dao.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        payment.DriverID,
		Type:          dao.TransactionTypeRelease,
		Amount:        payment.DriverAmount,
		Currency:      "ARS",
		BalanceBefore: wallet.AvailableBalance - payment.DriverAmount,
		BalanceAfter:  wallet.AvailableBalance,
		PaymentID:     &payment.PaymentUUID,
		BookingID:     &payment.BookingID,
		TripID:        &payment.TripID,
		Description:   "Fondos liberados automáticamente",
	}

	if err := s.walletRepo.CreateTransaction(transaction); err != nil {
		return err
	}

	// Update payment status
	payment.FundStatus = dao.FundStatusReleased
	payment.ReleasedAt = &now

	log.Info().
		Str("payment_uuid", payment.PaymentUUID).
		Int64("amount", payment.DriverAmount).
		Msg("Payment auto-released")

	return s.paymentRepo.Update(payment)
}
