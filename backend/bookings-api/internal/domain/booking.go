package domain

import "time"

// Request DTOs

type CreateBookingRequest struct {
	TripID        string `json:"trip_id" binding:"required"`
	SeatsReserved int    `json:"seats_reserved" binding:"required,min=1"`
}

// Response DTOs

type BookingDTO struct {
	ID                 string     `json:"id"`
	TripID             string     `json:"trip_id"`
	PassengerID        int64      `json:"passenger_id"`
	DriverID           int64      `json:"driver_id"`
	SeatsReserved      int        `json:"seats_reserved"`
	PricePerSeat       float64    `json:"price_per_seat"`
	TotalAmount        float64    `json:"total_amount"`
	Status             string     `json:"status"`
	PaymentStatus      string     `json:"payment_status"`
	ArrivedSafely      bool       `json:"arrived_safely"`
	ArrivalConfirmedAt *time.Time `json:"arrival_confirmed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type BookingListResponse struct {
	Bookings   []*BookingDTO  `json:"bookings"`
	Pagination PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Standard API Response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
