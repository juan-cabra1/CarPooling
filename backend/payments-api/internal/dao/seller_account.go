package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SellerAccountStatus represents the status of a seller account
type SellerAccountStatus string

const (
	// SellerStatusPending - OAuth flow started but not completed
	SellerStatusPending SellerAccountStatus = "pending"

	// SellerStatusActive - Successfully linked to Mercado Pago
	SellerStatusActive SellerAccountStatus = "active"

	// SellerStatusRevoked - Access was revoked by user or expired
	SellerStatusRevoked SellerAccountStatus = "revoked"

	// SellerStatusError - Error during OAuth flow
	SellerStatusError SellerAccountStatus = "error"
)

// SellerAccount represents a driver's linked Mercado Pago account
// This enables split payments where drivers receive their share directly
//
// OAuth Flow:
//  1. Driver requests linking -> we generate OAuth URL
//  2. Driver authorizes on Mercado Pago
//  3. MP redirects back with authorization code
//  4. We exchange code for access_token and refresh_token
//  5. Tokens are encrypted and stored here
//
// Token Refresh:
//   - access_token expires (typically 6 hours)
//   - We use refresh_token to get new access_token
//   - If refresh fails, status becomes "revoked"
type SellerAccount struct {
	// ID is the internal database primary key
	ID uint `gorm:"primaryKey;autoIncrement" json:"-"`

	// SellerUUID is the public identifier for this seller account
	SellerUUID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"seller_uuid"`

	// UserID is the UUID of the driver from users-api
	// One driver can have only one seller account
	UserID string `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_id"`

	// MPUserID is the Mercado Pago user ID returned after OAuth
	MPUserID string `gorm:"type:varchar(50);index" json:"mp_user_id,omitempty"`

	// MPAccessToken is the encrypted OAuth access token
	// Used for API calls on behalf of the seller
	// ENCRYPTED with AES-256
	MPAccessToken string `gorm:"type:text" json:"-"`

	// MPRefreshToken is the encrypted OAuth refresh token
	// Used to obtain new access tokens when they expire
	// ENCRYPTED with AES-256
	MPRefreshToken string `gorm:"type:text" json:"-"`

	// MPPublicKey is the seller's public key (for checkout)
	MPPublicKey string `gorm:"type:varchar(100)" json:"mp_public_key,omitempty"`

	// TokenExpiresAt is when the current access token expires
	TokenExpiresAt *time.Time `gorm:"index" json:"token_expires_at,omitempty"`

	// Status of the seller account
	Status SellerAccountStatus `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"`

	// ErrorMessage stores error details if status is "error"
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`

	// LastRefreshAt is when we last refreshed the OAuth token
	LastRefreshAt *time.Time `json:"last_refresh_at,omitempty"`

	// CreatedAt is when the record was created
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// UpdatedAt is when the record was last updated
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the custom table name
func (SellerAccount) TableName() string {
	return "seller_accounts"
}

// BeforeCreate generates UUID before inserting
func (s *SellerAccount) BeforeCreate(tx *gorm.DB) error {
	if s.SellerUUID == "" {
		s.SellerUUID = uuid.New().String()
	}
	return nil
}

// IsActive checks if the seller account is active and can receive payments
func (s *SellerAccount) IsActive() bool {
	return s.Status == SellerStatusActive
}

// IsTokenExpired checks if the access token has expired
func (s *SellerAccount) IsTokenExpired() bool {
	if s.TokenExpiresAt == nil {
		return true
	}
	// Consider expired 5 minutes before actual expiry for safety margin
	return time.Now().Add(5 * time.Minute).After(*s.TokenExpiresAt)
}

// NeedsRefresh checks if the token should be refreshed
func (s *SellerAccount) NeedsRefresh() bool {
	return s.IsActive() && s.IsTokenExpired()
}
