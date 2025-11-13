package domain

import "fmt"

// AppError represents a structured application error with a code and message
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError with optional details
func NewAppError(code, message string, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Predefined errors for bookings domain
var (
	// Booking errors
	ErrBookingNotFound = &AppError{
		Code:    "BOOKING_NOT_FOUND",
		Message: "Booking not found",
	}
	ErrDuplicateBooking = &AppError{
		Code:    "DUPLICATE_BOOKING",
		Message: "You already have a booking for this trip",
	}
	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Not authorized to perform this action",
	}
	ErrCannotCancelCompleted = &AppError{
		Code:    "CANNOT_CANCEL_COMPLETED",
		Message: "Cannot cancel a completed booking",
	}
	ErrBookingAlreadyCancelled = &AppError{
		Code:    "BOOKING_ALREADY_CANCELLED",
		Message: "Booking has already been cancelled",
	}

	// Trip validation errors
	ErrTripNotFound = &AppError{
		Code:    "TRIP_NOT_FOUND",
		Message: "Trip not found",
	}
	ErrTripNotPublished = &AppError{
		Code:    "TRIP_NOT_PUBLISHED",
		Message: "Trip is not published",
	}
	ErrInsufficientSeats = &AppError{
		Code:    "INSUFFICIENT_SEATS",
		Message: "Insufficient seats available",
	}
	ErrCannotBookOwnTrip = &AppError{
		Code:    "CANNOT_BOOK_OWN_TRIP",
		Message: "Cannot book your own trip",
	}

	// External service errors
	ErrTripsAPIUnavailable = &AppError{
		Code:    "TRIPS_API_UNAVAILABLE",
		Message: "Trips service is temporarily unavailable",
	}
)

// WithDetails returns a new AppError with additional details
func (e *AppError) WithDetails(details interface{}) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

// WithMessage returns a new AppError with a custom message
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: message,
		Details: e.Details,
	}
}

// Wrap wraps an error with context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
