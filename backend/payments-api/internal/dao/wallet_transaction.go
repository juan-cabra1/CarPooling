package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionType represents the type of wallet transaction
type TransactionType string

const (
	// TransactionTypeCredit - Money added to available balance (fund release)
	TransactionTypeCredit TransactionType = "credit"

	// TransactionTypeDebit - Money removed from available balance (withdrawal)
	TransactionTypeDebit TransactionType = "debit"

	// TransactionTypeHold - Money moved to held balance (payment received)
	TransactionTypeHold TransactionType = "hold"

	// TransactionTypeRelease - Money moved from held to available (trip confirmed)
	TransactionTypeRelease TransactionType = "release"

	// TransactionTypeRefund - Money deducted due to refund
	TransactionTypeRefund TransactionType = "refund"

	// TransactionTypeAdjustment - Manual adjustment by admin
	TransactionTypeAdjustment TransactionType = "adjustment"
)

// WalletTransaction represents a movement in a driver's wallet
//
// Every change to wallet balances creates a transaction record.
// This provides:
//   - Complete audit trail
//   - Balance verification (sum of transactions = current balance)
//   - Transaction history for drivers
type WalletTransaction struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// TransactionUUID is the public identifier for this transaction
	TransactionUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"transaction_uuid"`

	// WalletID is the internal ID of the wallet (foreign key)
	WalletID uint `gorm:"index;not null" json:"-"`

	// UserID is the UUID of the driver (denormalized for queries)
	UserID string `gorm:"type:varchar(36);index;not null" json:"user_id"`

	// Type of transaction
	Type TransactionType `gorm:"type:varchar(20);index;not null" json:"type"`

	// ============================================
	// AMOUNTS
	// ============================================

	// Amount of the transaction in centavos (always positive)
	Amount int64 `gorm:"not null" json:"amount"`

	// Currency is the 3-letter currency code
	Currency string `gorm:"type:varchar(3);not null;default:'ARS'" json:"currency"`

	// BalanceBefore is the available balance before this transaction
	BalanceBefore int64 `gorm:"not null" json:"balance_before"`

	// BalanceAfter is the available balance after this transaction
	BalanceAfter int64 `gorm:"not null" json:"balance_after"`

	// ============================================
	// REFERENCES
	// ============================================

	// PaymentID is the UUID of the related payment (if applicable)
	PaymentID *string `gorm:"type:varchar(36);index" json:"payment_id,omitempty"`

	// WithdrawalID is the UUID of the related withdrawal (if applicable)
	WithdrawalID *string `gorm:"type:varchar(36);index" json:"withdrawal_id,omitempty"`

	// BookingID is the UUID of the related booking (for reference)
	BookingID *string `gorm:"type:varchar(36)" json:"booking_id,omitempty"`

	// TripID is the UUID of the related trip (for reference)
	TripID *string `gorm:"type:varchar(36)" json:"trip_id,omitempty"`

	// ============================================
	// METADATA
	// ============================================

	// Description is a human-readable description of the transaction
	// Examples:
	//   - "Pago recibido por viaje Buenos Aires → Rosario"
	//   - "Retiro a cuenta Mercado Pago"
	//   - "Reembolso por cancelación de viaje"
	Description string `gorm:"type:text" json:"description"`

	// Metadata stores additional JSON data for the transaction
	// Can include passenger name, trip details, etc.
	Metadata string `gorm:"type:json" json:"metadata,omitempty"`

	// CreatedAt is when the transaction was created
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName specifies the custom table name
func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

// BeforeCreate generates UUID before inserting
func (wt *WalletTransaction) BeforeCreate(tx *gorm.DB) error {
	if wt.TransactionUUID == "" {
		wt.TransactionUUID = uuid.New().String()
	}
	return nil
}

// GetAmountInPesos returns the amount in pesos for display
func (wt *WalletTransaction) GetAmountInPesos() float64 {
	return float64(wt.Amount) / 100
}

// IsCredit returns true if this transaction added money
func (wt *WalletTransaction) IsCredit() bool {
	return wt.Type == TransactionTypeCredit || wt.Type == TransactionTypeRelease
}

// IsDebit returns true if this transaction removed money
func (wt *WalletTransaction) IsDebit() bool {
	return wt.Type == TransactionTypeDebit || wt.Type == TransactionTypeRefund
}

// GetSignedAmount returns positive for credits, negative for debits
func (wt *WalletTransaction) GetSignedAmount() int64 {
	if wt.IsDebit() {
		return -wt.Amount
	}
	return wt.Amount
}

// GetSignedAmountInPesos returns the signed amount in pesos
func (wt *WalletTransaction) GetSignedAmountInPesos() float64 {
	return float64(wt.GetSignedAmount()) / 100
}
