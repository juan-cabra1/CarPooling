package service

import (
	"time"

	"payments-api/internal/clients"
	"payments-api/internal/config"
	"payments-api/internal/dao"
	"payments-api/internal/domain"
	"payments-api/internal/repository"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// WithdrawalService handles withdrawal business logic
type WithdrawalService interface {
	RequestWithdrawal(userID string, amount int64) (*domain.WithdrawalResponse, error)
	GetWithdrawal(userID, withdrawalUUID string) (*domain.WithdrawalResponse, error)
	GetWithdrawals(userID string, page, pageSize int) (*domain.WithdrawalsListResponse, error)
	ProcessPendingWithdrawals() error
}

type withdrawalService struct {
	cfg            *config.Config
	withdrawalRepo repository.WithdrawalRepository
	walletRepo     repository.WalletRepository
	sellerRepo     repository.SellerRepository
	mpClient       clients.MercadoPagoClient
}

// NewWithdrawalService creates a new withdrawal service
func NewWithdrawalService(
	cfg *config.Config,
	withdrawalRepo repository.WithdrawalRepository,
	walletRepo repository.WalletRepository,
	sellerRepo repository.SellerRepository,
	mpClient clients.MercadoPagoClient,
) WithdrawalService {
	return &withdrawalService{
		cfg:            cfg,
		withdrawalRepo: withdrawalRepo,
		walletRepo:     walletRepo,
		sellerRepo:     sellerRepo,
		mpClient:       mpClient,
	}
}

// RequestWithdrawal creates a new withdrawal request
func (s *withdrawalService) RequestWithdrawal(userID string, amount int64) (*domain.WithdrawalResponse, error) {
	// Validate amount
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// Check for pending withdrawal
	hasPending, err := s.withdrawalRepo.HasPendingWithdrawal(userID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, domain.ErrWithdrawalInProgress
	}

	// Get wallet and verify balance
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	if !wallet.CanWithdraw(amount) {
		return nil, domain.ErrInsufficientBalance
	}

	// Check if seller account is active
	seller, err := s.sellerRepo.FindActiveByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSellerNotActive
		}
		return nil, err
	}

	// Deduct from available balance
	wallet.AvailableBalance -= amount
	now := time.Now()
	wallet.LastWithdrawalAt = &now

	if err := s.walletRepo.Update(wallet); err != nil {
		return nil, err
	}

	// Create withdrawal record
	withdrawal := &dao.Withdrawal{
		UserID:             userID,
		WalletID:           wallet.ID,
		Amount:             amount,
		Fee:                0, // No fee for now
		NetAmount:          amount,
		Currency:           "ARS",
		Status:             dao.WithdrawalStatusPending,
		DestinationAccount: seller.MPUserID,
	}

	if err := s.withdrawalRepo.Create(withdrawal); err != nil {
		// Rollback wallet balance
		wallet.AvailableBalance += amount
		s.walletRepo.Update(wallet)
		return nil, err
	}

	// Create transaction record
	transaction := &dao.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          dao.TransactionTypeDebit,
		Amount:        amount,
		Currency:      "ARS",
		BalanceBefore: wallet.AvailableBalance + amount,
		BalanceAfter:  wallet.AvailableBalance,
		WithdrawalID:  &withdrawal.WithdrawalUUID,
		Description:   "Retiro solicitado",
	}

	if err := s.walletRepo.CreateTransaction(transaction); err != nil {
		log.Error().Err(err).Msg("Failed to create wallet transaction for withdrawal")
	}

	log.Info().
		Str("withdrawal_uuid", withdrawal.WithdrawalUUID).
		Str("user_id", userID).
		Int64("amount", amount).
		Msg("Withdrawal requested")

	return s.toWithdrawalResponse(withdrawal), nil
}

// GetWithdrawal retrieves a withdrawal by UUID
func (s *withdrawalService) GetWithdrawal(userID, withdrawalUUID string) (*domain.WithdrawalResponse, error) {
	withdrawal, err := s.withdrawalRepo.FindByUUID(withdrawalUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWithdrawalNotFound
		}
		return nil, err
	}

	// Verify ownership
	if withdrawal.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return s.toWithdrawalResponse(withdrawal), nil
}

// GetWithdrawals retrieves withdrawals with pagination
func (s *withdrawalService) GetWithdrawals(userID string, page, pageSize int) (*domain.WithdrawalsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	withdrawals, total, err := s.withdrawalRepo.FindByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := &domain.WithdrawalsListResponse{
		Withdrawals: make([]domain.WithdrawalResponse, len(withdrawals)),
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}

	for i, w := range withdrawals {
		response.Withdrawals[i] = *s.toWithdrawalResponse(&w)
	}

	return response, nil
}

// ProcessPendingWithdrawals processes all pending withdrawals
// This should be called by a scheduler
func (s *withdrawalService) ProcessPendingWithdrawals() error {
	withdrawals, err := s.withdrawalRepo.FindPending()
	if err != nil {
		return err
	}

	for _, withdrawal := range withdrawals {
		if err := s.processWithdrawal(&withdrawal); err != nil {
			log.Error().
				Err(err).
				Str("withdrawal_uuid", withdrawal.WithdrawalUUID).
				Msg("Failed to process withdrawal")
		}
	}

	return nil
}

// processWithdrawal processes a single withdrawal
func (s *withdrawalService) processWithdrawal(withdrawal *dao.Withdrawal) error {
	// Update status to processing
	withdrawal.Status = dao.WithdrawalStatusProcessing
	now := time.Now()
	withdrawal.ProcessedAt = &now

	if err := s.withdrawalRepo.Update(withdrawal); err != nil {
		return err
	}

	// TODO: Call MP API to transfer funds
	// For now, simulate success
	withdrawal.Status = dao.WithdrawalStatusCompleted
	withdrawal.CompletedAt = &now

	// Update wallet total withdrawn
	wallet, err := s.walletRepo.FindByUserID(withdrawal.UserID)
	if err == nil {
		wallet.TotalWithdrawn += withdrawal.Amount
		s.walletRepo.Update(wallet)
	}

	return s.withdrawalRepo.Update(withdrawal)
}

// toWithdrawalResponse converts a withdrawal DAO to response
func (s *withdrawalService) toWithdrawalResponse(withdrawal *dao.Withdrawal) *domain.WithdrawalResponse {
	return &domain.WithdrawalResponse{
		WithdrawalUUID:     withdrawal.WithdrawalUUID,
		Amount:             withdrawal.Amount,
		Fee:                withdrawal.Fee,
		NetAmount:          withdrawal.NetAmount,
		Currency:           withdrawal.Currency,
		Status:             string(withdrawal.Status),
		DestinationAccount: withdrawal.DestinationAccount,
		ErrorMessage:       withdrawal.ErrorMessage,
		ProcessedAt:        withdrawal.ProcessedAt,
		CompletedAt:        withdrawal.CompletedAt,
		CreatedAt:          withdrawal.CreatedAt,
		AmountFormatted:    formatMoney(withdrawal.Amount),
	}
}
