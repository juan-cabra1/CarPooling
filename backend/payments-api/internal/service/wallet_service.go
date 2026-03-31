package service

import (
	"fmt"

	"payments-api/internal/dao"
	"payments-api/internal/domain"
	"payments-api/internal/repository"

	"gorm.io/gorm"
)

// WalletService handles wallet business logic
type WalletService interface {
	GetWallet(userID string) (*domain.WalletResponse, error)
	GetTransactions(userID string, page, pageSize int) (*domain.TransactionsListResponse, error)
	UpdateSettings(userID string, req *domain.WalletSettingsRequest) (*domain.WalletResponse, error)
	GetStats(userID string) (*domain.WalletStatsResponse, error)
}

type walletService struct {
	db         *gorm.DB
	walletRepo repository.WalletRepository
}

// NewWalletService creates a new wallet service
func NewWalletService(db *gorm.DB, walletRepo repository.WalletRepository) WalletService {
	return &walletService{
		db:         db,
		walletRepo: walletRepo,
	}
}

// GetWallet retrieves or creates a wallet for a user
func (s *walletService) GetWallet(userID string) (*domain.WalletResponse, error) {
	wallet, err := s.walletRepo.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}

	return s.toWalletResponse(wallet), nil
}

// GetTransactions retrieves wallet transactions with pagination
func (s *walletService) GetTransactions(userID string, page, pageSize int) (*domain.TransactionsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	transactions, total, err := s.walletRepo.FindTransactionsByUserID(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := &domain.TransactionsListResponse{
		Transactions: make([]domain.TransactionResponse, len(transactions)),
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}

	for i, tx := range transactions {
		response.Transactions[i] = s.toTransactionResponse(&tx)
	}

	return response, nil
}

// UpdateSettings updates wallet settings
func (s *walletService) UpdateSettings(userID string, req *domain.WalletSettingsRequest) (*domain.WalletResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	if req.AutoWithdrawEnabled != nil {
		wallet.AutoWithdrawEnabled = *req.AutoWithdrawEnabled
	}

	if req.AutoWithdrawThreshold != nil {
		if *req.AutoWithdrawThreshold < 0 {
			return nil, domain.ErrInvalidAmount
		}
		wallet.AutoWithdrawThreshold = *req.AutoWithdrawThreshold
	}

	if err := s.walletRepo.Update(wallet); err != nil {
		return nil, err
	}

	return s.toWalletResponse(wallet), nil
}

// GetStats retrieves wallet statistics
func (s *walletService) GetStats(userID string) (*domain.WalletStatsResponse, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	// TODO: Calculate monthly stats from transactions
	// For now, return basic stats
	stats := &domain.WalletStatsResponse{
		TotalEarnings:    wallet.TotalEarned,
		TotalTrips:       wallet.TotalTrips,
		TotalWithdrawals: wallet.TotalWithdrawn,
	}

	if wallet.TotalTrips > 0 {
		stats.AveragePerTrip = wallet.TotalEarned / int64(wallet.TotalTrips)
	}

	return stats, nil
}

// toWalletResponse converts a wallet DAO to response
func (s *walletService) toWalletResponse(wallet *dao.Wallet) *domain.WalletResponse {
	return &domain.WalletResponse{
		WalletUUID:                wallet.WalletUUID,
		UserID:                    wallet.UserID,
		AvailableBalance:          wallet.AvailableBalance,
		PendingBalance:            wallet.PendingBalance,
		HeldBalance:               wallet.HeldBalance,
		TotalBalance:              wallet.GetTotalBalance(),
		Currency:                  wallet.Currency,
		TotalEarned:               wallet.TotalEarned,
		TotalWithdrawn:            wallet.TotalWithdrawn,
		TotalTrips:                wallet.TotalTrips,
		AutoWithdrawEnabled:       wallet.AutoWithdrawEnabled,
		AutoWithdrawThreshold:     wallet.AutoWithdrawThreshold,
		AvailableBalanceFormatted: formatMoney(wallet.AvailableBalance),
		TotalEarnedFormatted:      formatMoney(wallet.TotalEarned),
	}
}

// toTransactionResponse converts a transaction DAO to response
func (s *walletService) toTransactionResponse(tx *dao.WalletTransaction) domain.TransactionResponse {
	signedAmount := tx.Amount
	if tx.IsDebit() {
		signedAmount = -signedAmount
	}

	return domain.TransactionResponse{
		TransactionUUID: tx.TransactionUUID,
		Type:            string(tx.Type),
		Amount:          tx.Amount,
		Currency:        tx.Currency,
		BalanceBefore:   tx.BalanceBefore,
		BalanceAfter:    tx.BalanceAfter,
		Description:     tx.Description,
		PaymentID:       tx.PaymentID,
		WithdrawalID:    tx.WithdrawalID,
		CreatedAt:       tx.CreatedAt,
		AmountFormatted: formatMoney(tx.Amount),
		SignedAmount:    signedAmount,
	}
}

// formatMoney formats centavos to a readable string
func formatMoney(centavos int64) string {
	pesos := float64(centavos) / 100
	return fmt.Sprintf("$%.2f", pesos)
}
