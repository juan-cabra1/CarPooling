package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WithdrawalStatus represents the status of a withdrawal request
type WithdrawalStatus string

const (
	// WithdrawalStatusPending - Request created, waiting to be processed
	WithdrawalStatusPending WithdrawalStatus = "pending"

	// WithdrawalStatusProcessing - Transfer initiated with Mercado Pago
	WithdrawalStatusProcessing WithdrawalStatus = "processing"

	// WithdrawalStatusCompleted - Transfer completed successfully
	WithdrawalStatusCompleted WithdrawalStatus = "completed"

	// WithdrawalStatusFailed - Transfer failed (will be retried or refunded to wallet)
	WithdrawalStatusFailed WithdrawalStatus = "failed"

	// WithdrawalStatusCancelled - Withdrawal cancelled by user or system
	WithdrawalStatusCancelled WithdrawalStatus = "cancelled"
)

// Withdrawal represents a request to withdraw funds from wallet to MP account
//
// Withdrawal Flow:
//  1. Driver requests withdrawal -> status "pending"
//  2. System validates balance and initiates MP transfer -> status "processing"
//  3. MP confirms transfer -> status "completed"
//  4. If transfer fails -> status "failed", funds returned to wallet
type Withdrawal struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// WithdrawalUUID is the public identifier for this withdrawal
	WithdrawalUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"withdrawal_uuid"`

	// UserID is the UUID of the driver requesting withdrawal
	UserID string `gorm:"type:varchar(36);index;not null" json:"user_id"`

	// WalletID is the internal ID of the wallet (foreign key)
	WalletID uint `gorm:"index;not null" json:"-"`

	// ============================================
	// AMOUNTS
	// ============================================

	// Amount is the withdrawal amount in centavos
	Amount int64 `gorm:"not null" json:"amount"`

	// Fee is any fee charged for the withdrawal (if applicable)
	// Currently 0, but may be added for instant withdrawals
	Fee int64 `gorm:"not null;default:0" json:"fee"`

	// NetAmount is the amount actually transferred (Amount - Fee)
	NetAmount int64 `gorm:"not null" json:"net_amount"`

	// Currency is the 3-letter currency code
	Currency string `gorm:"type:varchar(3);not null;default:'ARS'" json:"currency"`

	// ============================================
	// MERCADO PAGO REFERENCES
	// ============================================

	// MPTransferID is the Mercado Pago transfer/disbursement ID
	MPTransferID string `gorm:"type:varchar(100);index" json:"mp_transfer_id,omitempty"`

	// MPStatus is the raw status from Mercado Pago
	MPStatus string `gorm:"type:varchar(50)" json:"mp_status,omitempty"`

	// DestinationAccount contains info about where the funds went
	// Could be CVU, bank account, or MP wallet
	DestinationAccount string `gorm:"type:varchar(100)" json:"destination_account,omitempty"`

	// ============================================
	// STATUS
	// ============================================

	// Status of the withdrawal
	Status WithdrawalStatus `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"`

	// ErrorMessage stores error details if status is "failed"
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`

	// RetryCount tracks how many times we've retried this withdrawal
	RetryCount int `gorm:"not null;default:0" json:"retry_count"`

	// ============================================
	// TIMESTAMPS
	// ============================================

	// ProcessedAt is when the withdrawal started processing
	ProcessedAt *time.Time `json:"processed_at,omitempty"`

	// CompletedAt is when the withdrawal was completed
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// CreatedAt is when the record was created
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	// UpdatedAt is when the record was last updated
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the custom table name
func (Withdrawal) TableName() string {
	return "withdrawals"
}

// BeforeCreate generates UUID and calculates net amount before inserting
func (w *Withdrawal) BeforeCreate(tx *gorm.DB) error {
	if w.WithdrawalUUID == "" {
		w.WithdrawalUUID = uuid.New().String()
	}
	// Calculate net amount if not set
	if w.NetAmount == 0 {
		w.NetAmount = w.Amount - w.Fee
	}
	return nil
}

// IsCompleted checks if the withdrawal was successful
func (w *Withdrawal) IsCompleted() bool {
	return w.Status == WithdrawalStatusCompleted
}

// IsPending checks if the withdrawal is still pending
func (w *Withdrawal) IsPending() bool {
	return w.Status == WithdrawalStatusPending
}

// IsProcessing checks if the withdrawal is being processed
func (w *Withdrawal) IsProcessing() bool {
	return w.Status == WithdrawalStatusProcessing
}

// IsFailed checks if the withdrawal failed
func (w *Withdrawal) IsFailed() bool {
	return w.Status == WithdrawalStatusFailed
}

// CanBeRetried checks if a failed withdrawal can be retried
func (w *Withdrawal) CanBeRetried() bool {
	return w.Status == WithdrawalStatusFailed && w.RetryCount < 3
}

// GetAmountInPesos returns the amount in pesos for display
func (w *Withdrawal) GetAmountInPesos() float64 {
	return float64(w.Amount) / 100
}

// GetNetAmountInPesos returns the net amount in pesos for display
func (w *Withdrawal) GetNetAmountInPesos() float64 {
	return float64(w.NetAmount) / 100
}

// GetFeeInPesos returns the fee in pesos for display
func (w *Withdrawal) GetFeeInPesos() float64 {
	return float64(w.Fee) / 100
}
