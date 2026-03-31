package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"payments-api/internal/config"
	"payments-api/internal/domain"

	"github.com/rs/zerolog/log"
)

const (
	mpBaseURL     = "https://api.mercadopago.com"
	mpOAuthURL    = "https://auth.mercadopago.com"
	mpPreferenceURL = "/checkout/preferences"
	mpPaymentURL    = "/v1/payments"
	mpRefundURL     = "/v1/payments/%d/refunds"
	mpOAuthTokenURL = "/oauth/token"
)

// MercadoPagoClient handles communication with Mercado Pago API
type MercadoPagoClient interface {
	// Preferences
	CreatePreference(req *domain.MPPreferenceRequest) (*domain.MPPreferenceResponse, error)
	CreatePreferenceWithToken(req *domain.MPPreferenceRequest, accessToken string) (*domain.MPPreferenceResponse, error)

	// Payments
	GetPayment(paymentID string) (*domain.MPPaymentResponse, error)
	CreateRefund(paymentID int64, amount *float64) (*domain.MPRefundResponse, error)

	// OAuth
	GetOAuthURL(redirectURI, state string) string
	ExchangeCodeForTokens(code, redirectURI string) (*domain.MPOAuthTokenResponse, error)
	RefreshTokens(refreshToken string) (*domain.MPOAuthTokenResponse, error)
	GetUserInfo(accessToken string) (*domain.MPUserInfoResponse, error)
}

type mercadoPagoClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewMercadoPagoClient creates a new Mercado Pago client
func NewMercadoPagoClient(cfg *config.Config) MercadoPagoClient {
	return &mercadoPagoClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePreference creates a payment preference
func (c *mercadoPagoClient) CreatePreference(req *domain.MPPreferenceRequest) (*domain.MPPreferenceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preference request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", mpBaseURL+mpPreferenceURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, c.cfg.MPAccessToken)

	log.Debug().
		Str("url", mpBaseURL+mpPreferenceURL).
		Str("external_reference", req.ExternalReference).
		Msg("Creating MP preference")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("MP preference creation failed")
		return nil, fmt.Errorf("MP API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var prefResp domain.MPPreferenceResponse
	if err := json.Unmarshal(respBody, &prefResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Str("preference_id", prefResp.ID).
		Str("external_reference", req.ExternalReference).
		Msg("MP preference created successfully")

	return &prefResp, nil
}

// CreatePreferenceWithToken creates a payment preference using a specific access token (for marketplace)
func (c *mercadoPagoClient) CreatePreferenceWithToken(req *domain.MPPreferenceRequest, accessToken string) (*domain.MPPreferenceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preference request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", mpBaseURL+mpPreferenceURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, accessToken)

	log.Debug().
		Str("url", mpBaseURL+mpPreferenceURL).
		Str("external_reference", req.ExternalReference).
		Float64("marketplace_fee", req.MarketplaceFee).
		Msg("Creating MP marketplace preference")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("MP marketplace preference creation failed")
		return nil, fmt.Errorf("MP API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var prefResp domain.MPPreferenceResponse
	if err := json.Unmarshal(respBody, &prefResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Str("preference_id", prefResp.ID).
		Str("external_reference", req.ExternalReference).
		Float64("marketplace_fee", req.MarketplaceFee).
		Msg("MP marketplace preference created successfully")

	return &prefResp, nil
}

// GetPayment retrieves payment details from MP
func (c *mercadoPagoClient) GetPayment(paymentID string) (*domain.MPPaymentResponse, error) {
	url := fmt.Sprintf("%s%s/%s", mpBaseURL, mpPaymentURL, paymentID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, c.cfg.MPAccessToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrMPPaymentNotFound
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("MP get payment failed")
		return nil, fmt.Errorf("MP API error: status %d", resp.StatusCode)
	}

	var paymentResp domain.MPPaymentResponse
	if err := json.Unmarshal(respBody, &paymentResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &paymentResp, nil
}

// CreateRefund creates a refund for a payment
func (c *mercadoPagoClient) CreateRefund(paymentID int64, amount *float64) (*domain.MPRefundResponse, error) {
	url := fmt.Sprintf(mpBaseURL+mpRefundURL, paymentID)

	var body []byte
	var err error
	if amount != nil {
		req := domain.MPRefundRequest{Amount: *amount}
		body, err = json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal refund request: %w", err)
		}
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, c.cfg.MPAccessToken)

	log.Debug().
		Int64("payment_id", paymentID).
		Msg("Creating MP refund")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("MP refund creation failed")
		return nil, fmt.Errorf("MP API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var refundResp domain.MPRefundResponse
	if err := json.Unmarshal(respBody, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Int64("refund_id", refundResp.ID).
		Int64("payment_id", paymentID).
		Float64("amount", refundResp.Amount).
		Msg("MP refund created successfully")

	return &refundResp, nil
}

// GetOAuthURL returns the URL to redirect users for OAuth authorization
func (c *mercadoPagoClient) GetOAuthURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("client_id", c.cfg.MPClientID)
	params.Set("response_type", "code")
	params.Set("platform_id", "mp")
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)

	// Add scopes for marketplace functionality
	// read: read user data
	// offline_access: get refresh token
	params.Set("scope", "read offline_access")

	return fmt.Sprintf("%s/authorization?%s", mpOAuthURL, params.Encode())
}

// ExchangeCodeForTokens exchanges authorization code for access tokens
func (c *mercadoPagoClient) ExchangeCodeForTokens(code, redirectURI string) (*domain.MPOAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.cfg.MPClientID)
	data.Set("client_secret", c.cfg.MPClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	httpReq, err := http.NewRequest("POST", mpOAuthURL+mpOAuthTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	log.Debug().Msg("Exchanging OAuth code for tokens")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("OAuth token exchange failed")
		return nil, fmt.Errorf("OAuth error: status %d", resp.StatusCode)
	}

	var tokenResp domain.MPOAuthTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Int64("user_id", tokenResp.UserID).
		Msg("OAuth tokens obtained successfully")

	return &tokenResp, nil
}

// RefreshTokens refreshes expired access tokens
func (c *mercadoPagoClient) RefreshTokens(refreshToken string) (*domain.MPOAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.cfg.MPClientID)
	data.Set("client_secret", c.cfg.MPClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	httpReq, err := http.NewRequest("POST", mpOAuthURL+mpOAuthTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Accept", "application/json")

	log.Debug().Msg("Refreshing OAuth tokens")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("OAuth token refresh failed")
		return nil, fmt.Errorf("OAuth refresh error: status %d", resp.StatusCode)
	}

	var tokenResp domain.MPOAuthTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Int64("user_id", tokenResp.UserID).
		Msg("OAuth tokens refreshed successfully")

	return &tokenResp, nil
}

// GetUserInfo retrieves user information from MP using an access token
func (c *mercadoPagoClient) GetUserInfo(accessToken string) (*domain.MPUserInfoResponse, error) {
	url := fmt.Sprintf("%s/users/me", mpBaseURL)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Accept", "application/json")

	log.Debug().Msg("Getting MP user info")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("Get user info failed")
		return nil, fmt.Errorf("MP API error: status %d", resp.StatusCode)
	}

	var userInfo domain.MPUserInfoResponse
	if err := json.Unmarshal(respBody, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Info().
		Int64("user_id", userInfo.ID).
		Str("email", userInfo.Email).
		Msg("MP user info obtained")

	return &userInfo, nil
}

// setHeaders sets common headers for MP API requests
func (c *mercadoPagoClient) setHeaders(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("%d", time.Now().UnixNano()))
}
