package domain

import "time"

// WalletResponse is the wallet details response
type WalletResponse struct {
	WalletUUID        string  `json:"wallet_uuid"`
	UserID            string  `json:"user_id"`
	AvailableBalance  int64   `json:"available_balance"`
	PendingBalance    int64   `json:"pending_balance"`
	HeldBalance       int64   `json:"held_balance"`
	TotalBalance      int64   `json:"total_balance"`
	Currency          string  `json:"currency"`
	TotalEarned       int64   `json:"total_earned"`
	TotalWithdrawn    int64   `json:"total_withdrawn"`
	TotalTrips        int     `json:"total_trips"`
	AutoWithdrawEnabled   bool  `json:"auto_withdraw_enabled"`
	AutoWithdrawThreshold int64 `json:"auto_withdraw_threshold"`
	// Formatted for display
	AvailableBalanceFormatted string `json:"available_balance_formatted"`
	TotalEarnedFormatted      string `json:"total_earned_formatted"`
}

// WalletSettingsRequest is the request to update wallet settings
type WalletSettingsRequest struct {
	AutoWithdrawEnabled   *bool  `json:"auto_withdraw_enabled"`
	AutoWithdrawThreshold *int64 `json:"auto_withdraw_threshold"`
}

// TransactionResponse is a wallet transaction response
type TransactionResponse struct {
	TransactionUUID string    `json:"transaction_uuid"`
	Type            string    `json:"type"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	BalanceBefore   int64     `json:"balance_before"`
	BalanceAfter    int64     `json:"balance_after"`
	Description     string    `json:"description"`
	PaymentID       *string   `json:"payment_id,omitempty"`
	WithdrawalID    *string   `json:"withdrawal_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	// Formatted
	AmountFormatted string `json:"amount_formatted"`
	SignedAmount    int64  `json:"signed_amount"` // Positive for credit, negative for debit
}

// TransactionsListResponse is paginated list of transactions
type TransactionsListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// WalletStatsResponse contains wallet statistics
type WalletStatsResponse struct {
	// Current period (this month)
	CurrentMonthEarnings int64 `json:"current_month_earnings"`
	CurrentMonthTrips    int   `json:"current_month_trips"`

	// Previous period (last month)
	LastMonthEarnings int64 `json:"last_month_earnings"`
	LastMonthTrips    int   `json:"last_month_trips"`

	// Growth
	EarningsGrowthPercent float64 `json:"earnings_growth_percent"`

	// All time
	TotalEarnings    int64 `json:"total_earnings"`
	TotalTrips       int   `json:"total_trips"`
	TotalWithdrawals int64 `json:"total_withdrawals"`

	// Average
	AveragePerTrip int64 `json:"average_per_trip"`
}

// WithdrawalRequest is the request to withdraw funds
type WithdrawalRequest struct {
	Amount int64 `json:"amount" binding:"required,min=1"` // In centavos
}

// WithdrawalResponse is the withdrawal details response
type WithdrawalResponse struct {
	WithdrawalUUID     string     `json:"withdrawal_uuid"`
	Amount             int64      `json:"amount"`
	Fee                int64      `json:"fee"`
	NetAmount          int64      `json:"net_amount"`
	Currency           string     `json:"currency"`
	Status             string     `json:"status"`
	DestinationAccount string     `json:"destination_account,omitempty"`
	ErrorMessage       string     `json:"error_message,omitempty"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	// Formatted
	AmountFormatted string `json:"amount_formatted"`
}

// WithdrawalsListResponse is paginated list of withdrawals
type WithdrawalsListResponse struct {
	Withdrawals []WithdrawalResponse `json:"withdrawals"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	TotalPages  int                  `json:"total_pages"`
}
