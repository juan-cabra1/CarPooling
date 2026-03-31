package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	// PaymentStatusPending - Payment preference created, waiting for user to pay
	PaymentStatusPending PaymentStatus = "pending"

	// PaymentStatusApproved - Payment was approved by Mercado Pago
	PaymentStatusApproved PaymentStatus = "approved"

	// PaymentStatusRejected - Payment was rejected (insufficient funds, etc)
	PaymentStatusRejected PaymentStatus = "rejected"

	// PaymentStatusCancelled - Payment was cancelled before completion
	PaymentStatusCancelled PaymentStatus = "cancelled"

	// PaymentStatusRefunded - Payment was fully refunded
	PaymentStatusRefunded PaymentStatus = "refunded"

	// PaymentStatusPartialRefund - Payment was partially refunded
	PaymentStatusPartialRefund PaymentStatus = "partial_refund"

	// PaymentStatusInProcess - Payment is being processed
	PaymentStatusInProcess PaymentStatus = "in_process"

	// PaymentStatusChargedBack - Payment was disputed and charged back
	PaymentStatusChargedBack PaymentStatus = "charged_back"
)

// FundStatus represents the status of funds for this payment
type FundStatus string

const (
	// FundStatusHeld - Funds are held (escrow) pending trip completion
	FundStatusHeld FundStatus = "held"

	// FundStatusReleased - Funds have been released to driver's wallet
	FundStatusReleased FundStatus = "released"

	// FundStatusRefunded - Funds have been refunded to passenger
	FundStatusRefunded FundStatus = "refunded"
)

// Payment represents a payment for a booking
//
// Payment Flow:
//  1. Passenger books a trip -> Payment created with status "pending"
//  2. Passenger pays via Mercado Pago checkout
//  3. MP webhook notifies us -> status becomes "approved"
//  4. Booking is confirmed, funds are "held"
//  5. Trip completes -> after 2h or passenger confirmation, funds "released"
//  6. If cancelled >48h before: full refund, <48h: 50% refund
type Payment struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// PaymentUUID is the public identifier for this payment
	PaymentUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"payment_uuid"`

	// BookingID is the UUID of the booking from bookings-api
	BookingID string `gorm:"type:varchar(36);index;not null" json:"booking_id"`

	// TripID is the UUID of the trip from trips-api
	TripID string `gorm:"type:varchar(36);index;not null" json:"trip_id"`

	// PassengerID is the UUID of the passenger from users-api
	PassengerID string `gorm:"type:varchar(36);index;not null" json:"passenger_id"`

	// DriverID is the UUID of the driver from users-api
	DriverID string `gorm:"type:varchar(36);index;not null" json:"driver_id"`

	// ============================================
	// AMOUNTS (all in smallest currency unit, e.g., centavos)
	// ============================================

	// Amount is the total payment amount in centavos
	// Example: $100.00 = 10000 centavos
	Amount int64 `gorm:"not null" json:"amount"`

	// Currency is the 3-letter currency code (ISO 4217)
	// Default: "ARS" (Argentine Peso)
	Currency string `gorm:"type:varchar(3);not null;default:'ARS'" json:"currency"`

	// PlatformFee is the amount taken by the platform (15% default)
	// Stored in centavos
	PlatformFee int64 `gorm:"not null" json:"platform_fee"`

	// DriverAmount is the amount that goes to the driver (85% default)
	// Stored in centavos
	DriverAmount int64 `gorm:"not null" json:"driver_amount"`

	// SeatsCount is the number of seats in this payment
	SeatsCount int `gorm:"not null;default:1" json:"seats_count"`

	// PricePerSeat is the price per seat in centavos
	PricePerSeat int64 `gorm:"not null" json:"price_per_seat"`

	// ============================================
	// MERCADO PAGO REFERENCES
	// ============================================

	// MPPreferenceID is the Mercado Pago preference ID
	// Used to track the checkout session
	MPPreferenceID string `gorm:"type:varchar(100);index" json:"mp_preference_id,omitempty"`

	// MPPaymentID is the Mercado Pago payment ID
	// Set when payment is completed
	MPPaymentID string `gorm:"type:varchar(100);index" json:"mp_payment_id,omitempty"`

	// MPMerchantOrderID is the merchant order ID from MP
	MPMerchantOrderID string `gorm:"type:varchar(100)" json:"mp_merchant_order_id,omitempty"`

	// MPStatus is the raw status from Mercado Pago
	MPStatus string `gorm:"type:varchar(50)" json:"mp_status,omitempty"`

	// MPStatusDetail provides additional status details from MP
	MPStatusDetail string `gorm:"type:varchar(100)" json:"mp_status_detail,omitempty"`

	// ============================================
	// STATUS TRACKING
	// ============================================

	// Status of the payment
	Status PaymentStatus `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"`

	// FundStatus indicates what happened to the funds
	FundStatus FundStatus `gorm:"type:varchar(20);index;not null;default:'held'" json:"fund_status"`

	// ============================================
	// REFUND TRACKING
	// ============================================

	// RefundedAmount is the total amount refunded (in centavos)
	RefundedAmount int64 `gorm:"default:0" json:"refunded_amount"`

	// RefundReason explains why the refund was issued
	RefundReason string `gorm:"type:text" json:"refund_reason,omitempty"`

	// RefundedAt is when the refund was processed
	RefundedAt *time.Time `json:"refunded_at,omitempty"`

	// ============================================
	// TIMESTAMPS
	// ============================================

	// ApprovedAt is when the payment was approved by MP
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// ReleasedAt is when funds were released to driver's wallet
	ReleasedAt *time.Time `json:"released_at,omitempty"`

	// ExpiresAt is when the payment preference expires
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// CreatedAt is when the record was created
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// UpdatedAt is when the record was last updated
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the custom table name
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate generates UUID before inserting
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.PaymentUUID == "" {
		p.PaymentUUID = uuid.New().String()
	}
	return nil
}

// IsCompleted checks if payment was successfully completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusApproved
}

// IsPending checks if payment is still pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusInProcess
}

// CanBeRefunded checks if this payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusApproved && p.FundStatus == FundStatusHeld
}

// CanBeReleased checks if funds can be released to driver
func (p *Payment) CanBeReleased() bool {
	return p.Status == PaymentStatusApproved && p.FundStatus == FundStatusHeld
}

// GetAmountInPesos returns the amount in pesos (for display)
func (p *Payment) GetAmountInPesos() float64 {
	return float64(p.Amount) / 100
}

// GetDriverAmountInPesos returns the driver amount in pesos (for display)
func (p *Payment) GetDriverAmountInPesos() float64 {
	return float64(p.DriverAmount) / 100
}

// GetPlatformFeeInPesos returns the platform fee in pesos (for display)
func (p *Payment) GetPlatformFeeInPesos() float64 {
	return float64(p.PlatformFee) / 100
}
