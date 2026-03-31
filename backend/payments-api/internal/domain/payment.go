package domain

import "time"

// CreatePreferenceRequest is the request to create a payment preference
type CreatePreferenceRequest struct {
	BookingID   string `json:"booking_id" binding:"required"`
	TripID      string `json:"trip_id" binding:"required"`
	DriverID    string `json:"driver_id" binding:"required"`
	SeatsCount  int    `json:"seats_count" binding:"required,min=1"`
	PricePerSeat int64  `json:"price_per_seat" binding:"required,min=1"` // In centavos

	// Trip details for preference description
	Origin      string    `json:"origin" binding:"required"`
	Destination string    `json:"destination" binding:"required"`
	DepartureAt time.Time `json:"departure_at" binding:"required"`

	// Payer info
	PayerEmail string `json:"payer_email"`
	PayerName  string `json:"payer_name"`

	// ReturnBaseURL es la URL base a la que MP redirige después del pago.
	// Mobile: "carpooling:/" → deep link que abre la app
	// Web:    "https://tu-app.railway.app" → URL del navegador
	// Si está vacío, se usa el default del servidor (deep link mobile).
	ReturnBaseURL string `json:"return_base_url"`
}

// CreatePreferenceResponse is the response after creating a preference
type CreatePreferenceResponse struct {
	PaymentUUID    string `json:"payment_uuid"`
	PreferenceID   string `json:"preference_id"`
	InitPoint      string `json:"init_point"`       // URL for production checkout
	SandboxInitPoint string `json:"sandbox_init_point"` // URL for sandbox checkout
	Amount         int64  `json:"amount"`           // Total in centavos
	PlatformFee    int64  `json:"platform_fee"`     // Platform fee in centavos
	DriverAmount   int64  `json:"driver_amount"`    // Driver amount in centavos
	ExpiresAt      time.Time `json:"expires_at"`
}

// PaymentResponse is the payment details response
type PaymentResponse struct {
	PaymentUUID     string     `json:"payment_uuid"`
	BookingID       string     `json:"booking_id"`
	TripID          string     `json:"trip_id"`
	PassengerID     string     `json:"passenger_id"`
	DriverID        string     `json:"driver_id"`
	Amount          int64      `json:"amount"`
	Currency        string     `json:"currency"`
	PlatformFee     int64      `json:"platform_fee"`
	DriverAmount    int64      `json:"driver_amount"`
	SeatsCount      int        `json:"seats_count"`
	Status          string     `json:"status"`
	FundStatus      string     `json:"fund_status"`
	MPPaymentID     string     `json:"mp_payment_id,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	ReleasedAt      *time.Time `json:"released_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// ReleasePaymentRequest is the request to release funds
type ReleasePaymentRequest struct {
	PaymentUUID string `json:"payment_uuid" binding:"required"`
	// Optional: rating and comment from passenger
	Rating  *int   `json:"rating,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// RefundRequest is the request to refund a payment
type RefundRequest struct {
	PaymentUUID string  `json:"payment_uuid" binding:"required"`
	Reason      string  `json:"reason" binding:"required"`
	Percentage  float64 `json:"percentage"` // 0-100, default 100 for full refund
}
