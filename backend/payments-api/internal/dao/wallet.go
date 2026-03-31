package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wallet represents a driver's wallet for holding earnings
//
// Balance Types:
//   - AvailableBalance: Funds that can be withdrawn immediately
//   - PendingBalance: Funds from trips not yet completed (still in progress)
//   - HeldBalance: Funds held in escrow post-trip, waiting for release (2h timer)
//
// Flow:
//  1. Payment approved -> funds go to HeldBalance
//  2. Trip completed + 2h OR passenger confirms -> funds move to AvailableBalance
//  3. Driver requests withdrawal -> funds deducted from AvailableBalance
type Wallet struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// WalletUUID is the public identifier for this wallet
	WalletUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"wallet_uuid"`

	// UserID is the UUID of the driver from users-api
	// One driver has exactly one wallet
	UserID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"`

	// ============================================
	// BALANCES (all in smallest currency unit, e.g., centavos)
	// ============================================

	// AvailableBalance is the amount ready for withdrawal
	// Stored in centavos (e.g., $100.00 = 10000)
	AvailableBalance int64 `gorm:"not null;default:0" json:"available_balance"`

	// PendingBalance is the amount from trips in progress
	// These funds haven't been captured yet
	PendingBalance int64 `gorm:"not null;default:0" json:"pending_balance"`

	// HeldBalance is the amount held in escrow
	// Waiting for trip confirmation or 2h timer
	HeldBalance int64 `gorm:"not null;default:0" json:"held_balance"`

	// Currency is the 3-letter currency code (ISO 4217)
	Currency string `gorm:"type:varchar(3);not null;default:'ARS'" json:"currency"`

	// ============================================
	// STATISTICS
	// ============================================

	// TotalEarned is the lifetime total earnings (for stats)
	TotalEarned int64 `gorm:"not null;default:0" json:"total_earned"`

	// TotalWithdrawn is the lifetime total withdrawn (for stats)
	TotalWithdrawn int64 `gorm:"not null;default:0" json:"total_withdrawn"`

	// TotalTrips is the number of completed trips
	TotalTrips int `gorm:"not null;default:0" json:"total_trips"`

	// ============================================
	// AUTO-WITHDRAW SETTINGS
	// ============================================

	// AutoWithdrawEnabled enables automatic withdrawals
	AutoWithdrawEnabled bool `gorm:"not null;default:false" json:"auto_withdraw_enabled"`

	// AutoWithdrawThreshold is the minimum balance to trigger auto-withdrawal
	// Stored in centavos
	AutoWithdrawThreshold int64 `gorm:"not null;default:0" json:"auto_withdraw_threshold"`

	// ============================================
	// TIMESTAMPS
	// ============================================

	// LastWithdrawalAt is when the last withdrawal was made
	LastWithdrawalAt *time.Time `json:"last_withdrawal_at,omitempty"`

	// LastDepositAt is when the last deposit was made
	LastDepositAt *time.Time `json:"last_deposit_at,omitempty"`

	// CreatedAt is when the record was created
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// UpdatedAt is when the record was last updated
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the custom table name
func (Wallet) TableName() string {
	return "wallets"
}

// BeforeCreate generates UUID before inserting
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.WalletUUID == "" {
		w.WalletUUID = uuid.New().String()
	}
	return nil
}

// GetTotalBalance returns the sum of all balances
func (w *Wallet) GetTotalBalance() int64 {
	return w.AvailableBalance + w.PendingBalance + w.HeldBalance
}

// GetAvailableBalanceInPesos returns available balance in pesos
func (w *Wallet) GetAvailableBalanceInPesos() float64 {
	return float64(w.AvailableBalance) / 100
}

// GetPendingBalanceInPesos returns pending balance in pesos
func (w *Wallet) GetPendingBalanceInPesos() float64 {
	return float64(w.PendingBalance) / 100
}

// GetHeldBalanceInPesos returns held balance in pesos
func (w *Wallet) GetHeldBalanceInPesos() float64 {
	return float64(w.HeldBalance) / 100
}

// GetTotalBalanceInPesos returns total balance in pesos
func (w *Wallet) GetTotalBalanceInPesos() float64 {
	return float64(w.GetTotalBalance()) / 100
}

// GetTotalEarnedInPesos returns lifetime earnings in pesos
func (w *Wallet) GetTotalEarnedInPesos() float64 {
	return float64(w.TotalEarned) / 100
}

// GetTotalWithdrawnInPesos returns lifetime withdrawals in pesos
func (w *Wallet) GetTotalWithdrawnInPesos() float64 {
	return float64(w.TotalWithdrawn) / 100
}

// CanWithdraw checks if the wallet has enough available balance
func (w *Wallet) CanWithdraw(amount int64) bool {
	return w.AvailableBalance >= amount && amount > 0
}

// ShouldAutoWithdraw checks if auto-withdrawal should be triggered
func (w *Wallet) ShouldAutoWithdraw() bool {
	return w.AutoWithdrawEnabled &&
		w.AvailableBalance >= w.AutoWithdrawThreshold &&
		w.AutoWithdrawThreshold > 0
}
