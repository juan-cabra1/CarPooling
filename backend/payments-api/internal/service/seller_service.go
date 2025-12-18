package service

import (
	"fmt"
	"time"

	"payments-api/internal/clients"
	"payments-api/internal/config"
	"payments-api/internal/dao"
	"payments-api/internal/domain"
	"payments-api/internal/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// SellerService handles seller account business logic
type SellerService interface {
	GetOAuthURL(userID, redirectURI string) (string, error)
	ProcessOAuthCallback(userID, code, redirectURI string) error
	LinkManualToken(userID, accessToken string) error
	GetStatus(userID string) (*SellerStatusResponse, error)
	Disconnect(userID string) error
	RefreshExpiredTokens() error
}

// SellerStatusResponse represents the seller status response
type SellerStatusResponse struct {
	IsLinked      bool       `json:"is_linked"`
	Status        string     `json:"status"`
	MPUserID      string     `json:"mp_user_id,omitempty"`
	LinkedAt      *time.Time `json:"linked_at,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
}

type sellerService struct {
	cfg        *config.Config
	sellerRepo repository.SellerRepository
	mpClient   clients.MercadoPagoClient
}

// NewSellerService creates a new seller service
func NewSellerService(
	cfg *config.Config,
	sellerRepo repository.SellerRepository,
	mpClient clients.MercadoPagoClient,
) SellerService {
	return &sellerService{
		cfg:        cfg,
		sellerRepo: sellerRepo,
		mpClient:   mpClient,
	}
}

// GetOAuthURL returns the URL to redirect user for MP OAuth
func (s *sellerService) GetOAuthURL(userID, redirectURI string) (string, error) {
	// Check if already linked
	existing, err := s.sellerRepo.FindActiveByUserID(userID)
	if err == nil && existing != nil {
		return "", domain.ErrSellerAlreadyLinked
	}

	// Generate state for CSRF protection
	state := fmt.Sprintf("%s:%s", userID, uuid.New().String())

	// Create or update pending seller record
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		seller = &dao.SellerAccount{
			UserID: userID,
			Status: dao.SellerStatusPending,
		}
		if err := s.sellerRepo.Create(seller); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
		seller.Status = dao.SellerStatusPending
		seller.ErrorMessage = ""
		if err := s.sellerRepo.Update(seller); err != nil {
			return "", err
		}
	}

	return s.mpClient.GetOAuthURL(redirectURI, state), nil
}

// ProcessOAuthCallback processes the OAuth callback from MP
func (s *sellerService) ProcessOAuthCallback(userID, code, redirectURI string) error {
	// Exchange code for tokens
	tokens, err := s.mpClient.ExchangeCodeForTokens(code, redirectURI)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to exchange OAuth code")

		// Update seller with error
		seller, _ := s.sellerRepo.FindByUserID(userID)
		if seller != nil {
			seller.Status = dao.SellerStatusError
			seller.ErrorMessage = "Error al vincular cuenta de Mercado Pago"
			s.sellerRepo.Update(seller)
		}

		return domain.ErrOAuthFailed
	}

	// Get or create seller record
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		seller = &dao.SellerAccount{
			UserID: userID,
		}
	} else if err != nil {
		return err
	}

	// Update with token info
	seller.MPUserID = fmt.Sprintf("%d", tokens.UserID)
	seller.MPAccessToken = tokens.AccessToken  // TODO: Encrypt
	seller.MPRefreshToken = tokens.RefreshToken // TODO: Encrypt
	seller.MPPublicKey = tokens.PublicKey
	seller.Status = dao.SellerStatusActive
	seller.ErrorMessage = ""

	// Calculate token expiration
	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	seller.TokenExpiresAt = &expiresAt
	now := time.Now()
	seller.LastRefreshAt = &now

	if seller.ID == 0 {
		err = s.sellerRepo.Create(seller)
	} else {
		err = s.sellerRepo.Update(seller)
	}

	if err != nil {
		return err
	}

	log.Info().
		Str("user_id", userID).
		Str("mp_user_id", seller.MPUserID).
		Msg("Seller account linked successfully")

	return nil
}

// LinkManualToken links a seller account using a manually provided access token
func (s *sellerService) LinkManualToken(userID, accessToken string) error {
	// Check if already linked
	existing, err := s.sellerRepo.FindActiveByUserID(userID)
	if err == nil && existing != nil {
		return domain.ErrSellerAlreadyLinked
	}

	// Validate the token by making a test API call to get user info
	userInfo, err := s.mpClient.GetUserInfo(accessToken)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to validate MP token")
		return domain.ErrInvalidToken
	}

	// Get or create seller record
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		seller = &dao.SellerAccount{
			UserID: userID,
		}
	} else if err != nil {
		return err
	}

	// Update with token info
	seller.MPUserID = fmt.Sprintf("%d", userInfo.ID)
	seller.MPAccessToken = accessToken  // TODO: Encrypt
	seller.MPRefreshToken = ""          // No refresh token in manual mode
	seller.MPPublicKey = ""             // Not needed for manual mode
	seller.Status = dao.SellerStatusActive
	seller.ErrorMessage = ""
	seller.TokenExpiresAt = nil // Manual tokens don't expire automatically
	now := time.Now()
	seller.LastRefreshAt = &now

	if seller.ID == 0 {
		err = s.sellerRepo.Create(seller)
	} else {
		err = s.sellerRepo.Update(seller)
	}

	if err != nil {
		return err
	}

	log.Info().
		Str("user_id", userID).
		Str("mp_user_id", seller.MPUserID).
		Msg("Seller account linked manually")

	return nil
}

// GetStatus returns the seller account status
func (s *sellerService) GetStatus(userID string) (*SellerStatusResponse, error) {
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		return &SellerStatusResponse{
			IsLinked: false,
			Status:   "not_linked",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &SellerStatusResponse{
		IsLinked:     seller.IsActive(),
		Status:       string(seller.Status),
		MPUserID:     seller.MPUserID,
		LinkedAt:     &seller.CreatedAt,
		ErrorMessage: seller.ErrorMessage,
	}, nil
}

// Disconnect removes the MP account link
func (s *sellerService) Disconnect(userID string) error {
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ErrSellerNotFound
		}
		return err
	}

	seller.Status = dao.SellerStatusRevoked
	seller.MPAccessToken = ""
	seller.MPRefreshToken = ""

	return s.sellerRepo.Update(seller)
}

// RefreshExpiredTokens refreshes tokens that are about to expire
// This should be called by a scheduler
func (s *sellerService) RefreshExpiredTokens() error {
	sellers, err := s.sellerRepo.FindNeedingTokenRefresh()
	if err != nil {
		return err
	}

	for _, seller := range sellers {
		if err := s.refreshSellerToken(&seller); err != nil {
			log.Error().
				Err(err).
				Str("seller_uuid", seller.SellerUUID).
				Msg("Failed to refresh seller token")
		}
	}

	return nil
}

// refreshSellerToken refreshes a single seller's token
func (s *sellerService) refreshSellerToken(seller *dao.SellerAccount) error {
	tokens, err := s.mpClient.RefreshTokens(seller.MPRefreshToken)
	if err != nil {
		seller.Status = dao.SellerStatusRevoked
		seller.ErrorMessage = "Token refresh failed - please reconnect"
		return s.sellerRepo.Update(seller)
	}

	seller.MPAccessToken = tokens.AccessToken
	seller.MPRefreshToken = tokens.RefreshToken

	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	seller.TokenExpiresAt = &expiresAt
	now := time.Now()
	seller.LastRefreshAt = &now

	log.Info().
		Str("seller_uuid", seller.SellerUUID).
		Msg("Seller token refreshed")

	return s.sellerRepo.Update(seller)
}
