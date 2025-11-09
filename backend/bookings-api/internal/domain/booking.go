package domain

import (
	"bookings-api/internal/dao"
	"time"
)

// CreateBookingRequest represents the request to create a new booking
type CreateBookingRequest struct {
	TripID        string `json:"trip_id" binding:"required"`
	PassengerID   int64  `json:"passenger_id" binding:"required"`
	SeatsReserved int    `json:"seats_reserved" binding:"required,min=1"`
}

// BookingResponse represents a booking in API responses
type BookingResponse struct {
	ID                 string     `json:"id"`
	TripID             string     `json:"trip_id"`
	PassengerID        int64      `json:"passenger_id"`
	SeatsRequested     int        `json:"seats_requested"`
	TotalPrice         float64    `json:"total_price"`
	Status             string     `json:"status"`
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// CancelBookingRequest represents the request to cancel a booking
type CancelBookingRequest struct {
	Reason string `json:"reason"`
}

// BookingListResponse represents a paginated list of bookings
type BookingListResponse struct {
	Bookings   []BookingResponse `json:"bookings"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// Booking status constants (mirror DAO constants for clarity)
const (
	BookingStatusPending   = dao.BookingStatusPending
	BookingStatusConfirmed = dao.BookingStatusConfirmed
	BookingStatusCancelled = dao.BookingStatusCancelled
	BookingStatusCompleted = dao.BookingStatusCompleted
	BookingStatusFailed    = dao.BookingStatusFailed
)

// ToBookingResponse converts a DAO Booking to a BookingResponse DTO
func ToBookingResponse(b *dao.Booking) *BookingResponse {
	if b == nil {
		return nil
	}

	return &BookingResponse{
		ID:                 b.BookingUUID,
		TripID:             b.TripID,
		PassengerID:        b.PassengerID,
		SeatsRequested:     b.SeatsRequested,
		TotalPrice:         b.TotalPrice,
		Status:             b.Status,
		CancelledAt:        b.CancelledAt,
		CancellationReason: b.CancellationReason,
		CreatedAt:          b.CreatedAt,
		UpdatedAt:          b.UpdatedAt,
	}
}

// ToBookingResponseList converts a slice of DAO Bookings to BookingResponse DTOs
func ToBookingResponseList(bookings []dao.Booking) []BookingResponse {
	responses := make([]BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		responses = append(responses, *ToBookingResponse(&booking))
	}
	return responses
}

// CalculateTotalPages calculates the total number of pages for pagination
func CalculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 0
	}
	pages := int(total) / limit
	if int(total)%limit > 0 {
		pages++
	}
	return pages
}

// NewBookingListResponse creates a paginated booking list response
func NewBookingListResponse(bookings []dao.Booking, total int64, page, limit int) *BookingListResponse {
	return &BookingListResponse{
		Bookings:   ToBookingResponseList(bookings),
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: CalculateTotalPages(total, limit),
	}
}
