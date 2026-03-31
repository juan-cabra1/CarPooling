package domain

import "time"

// MPPreferenceItem represents an item in the preference
type MPPreferenceItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PictureURL  string  `json:"picture_url,omitempty"`
	CategoryID  string  `json:"category_id,omitempty"`
	Quantity    int     `json:"quantity"`
	CurrencyID  string  `json:"currency_id"`
	UnitPrice   float64 `json:"unit_price"`
}

// MPPayer represents the payer information
type MPPayer struct {
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
	Email   string `json:"email"`
	Phone   *MPPhone `json:"phone,omitempty"`
	Address *MPAddress `json:"address,omitempty"`
}

// MPPhone represents phone information
type MPPhone struct {
	AreaCode string `json:"area_code,omitempty"`
	Number   string `json:"number,omitempty"`
}

// MPAddress represents address information
type MPAddress struct {
	ZipCode    string `json:"zip_code,omitempty"`
	StreetName string `json:"street_name,omitempty"`
	StreetNumber int  `json:"street_number,omitempty"`
}

// MPBackURLs represents the callback URLs
type MPBackURLs struct {
	Success string `json:"success"`
	Pending string `json:"pending"`
	Failure string `json:"failure"`
}

// MPPreferenceRequest is the request to create a preference in MP
type MPPreferenceRequest struct {
	Items             []MPPreferenceItem `json:"items"`
	Payer             *MPPayer           `json:"payer,omitempty"`
	BackURLs          *MPBackURLs        `json:"back_urls,omitempty"`
	AutoReturn        string             `json:"auto_return,omitempty"` // "approved", "all"
	ExternalReference string             `json:"external_reference"`
	NotificationURL   string             `json:"notification_url,omitempty"`
	ExpirationDateFrom *time.Time        `json:"expiration_date_from,omitempty"`
	ExpirationDateTo   *time.Time        `json:"expiration_date_to,omitempty"`
	Expires           bool               `json:"expires,omitempty"`

	// Marketplace fields for split payments
	Marketplace        string  `json:"marketplace,omitempty"`
	MarketplaceFee     float64 `json:"marketplace_fee,omitempty"`

	// Additional info
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
}

// MPPreferenceResponse is the response from creating a preference
type MPPreferenceResponse struct {
	ID               string    `json:"id"`
	InitPoint        string    `json:"init_point"`
	SandboxInitPoint string    `json:"sandbox_init_point"`
	DateCreated      time.Time `json:"date_created"`
	ExpirationDateTo *time.Time `json:"expiration_date_to,omitempty"`
	ExternalReference string   `json:"external_reference"`
	CollectorID      int64     `json:"collector_id"`
}

// MPWebhookPayload is the webhook notification from MP
type MPWebhookPayload struct {
	ID          int64  `json:"id"`
	LiveMode    bool   `json:"live_mode"`
	Type        string `json:"type"` // "payment", "merchant_order", etc.
	DateCreated string `json:"date_created"`
	UserID      int64  `json:"user_id"`
	APIVersion  string `json:"api_version"`
	Action      string `json:"action"` // "payment.created", "payment.updated"
	Data        struct {
		ID string `json:"id"`
	} `json:"data"`
}

// MPPaymentResponse is the payment details from MP API
type MPPaymentResponse struct {
	ID                    int64     `json:"id"`
	DateCreated           time.Time `json:"date_created"`
	DateApproved          *time.Time `json:"date_approved,omitempty"`
	DateLastUpdated       time.Time `json:"date_last_updated"`
	MoneyReleaseDate      *time.Time `json:"money_release_date,omitempty"`
	Status                string    `json:"status"` // approved, pending, rejected, etc.
	StatusDetail          string    `json:"status_detail"`
	TransactionAmount     float64   `json:"transaction_amount"`
	CurrencyID            string    `json:"currency_id"`
	Description           string    `json:"description"`
	ExternalReference     string    `json:"external_reference"`
	PayerID               int64     `json:"payer_id,omitempty"`
	PaymentMethodID       string    `json:"payment_method_id"`
	PaymentTypeID         string    `json:"payment_type_id"`
	MerchantOrderID       int64     `json:"merchant_order_id,omitempty"`

	// Fee info
	FeeDetails []struct {
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
		FeePayer string  `json:"fee_payer"`
	} `json:"fee_details,omitempty"`

	// Transaction details
	TransactionDetails struct {
		NetReceivedAmount  float64 `json:"net_received_amount"`
		TotalPaidAmount    float64 `json:"total_paid_amount"`
		OverpaidAmount     float64 `json:"overpaid_amount"`
		InstallmentAmount  float64 `json:"installment_amount"`
	} `json:"transaction_details,omitempty"`
}

// MPRefundRequest is the request to create a refund
type MPRefundRequest struct {
	Amount float64 `json:"amount,omitempty"` // Optional: for partial refunds
}

// MPRefundResponse is the response from creating a refund
type MPRefundResponse struct {
	ID                    int64     `json:"id"`
	PaymentID             int64     `json:"payment_id"`
	Amount                float64   `json:"amount"`
	Status                string    `json:"status"`
	DateCreated           time.Time `json:"date_created"`
	Source                struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"source"`
}

// MPOAuthTokenRequest is the request to exchange code for tokens
type MPOAuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"` // "authorization_code"
}

// MPOAuthTokenResponse is the response with OAuth tokens
type MPOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // Seconds
	Scope        string `json:"scope"`
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	PublicKey    string `json:"public_key"`
	LiveMode     bool   `json:"live_mode"`
}

// MPOAuthRefreshRequest is the request to refresh tokens
type MPOAuthRefreshRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"` // "refresh_token"
}

// MPUserInfoResponse is the response from getting user info
type MPUserInfoResponse struct {
	ID          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	SiteID      string `json:"site_id"`
	CountryID   string `json:"country_id"`
}
