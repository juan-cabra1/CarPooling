package domain

import "fmt"

// AppError represents a structured application error
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

// Common errors for HTTP clients and external services
var (
	// HTTP Client Errors
	ErrTripNotFound = &AppError{
		Code:    "TRIP_NOT_FOUND",
		Message: "Trip not found in trips-api",
	}

	ErrUserNotFound = &AppError{
		Code:    "USER_NOT_FOUND",
		Message: "User not found in users-api",
	}

	ErrServiceUnavailable = &AppError{
		Code:    "SERVICE_UNAVAILABLE",
		Message: "External service temporarily unavailable",
	}

	ErrTimeout = &AppError{
		Code:    "REQUEST_TIMEOUT",
		Message: "Request to external service timed out",
	}

	ErrInvalidResponse = &AppError{
		Code:    "INVALID_RESPONSE",
		Message: "Received invalid response from external service",
	}

	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized access to external service",
	}

	// Repository Errors
	ErrSearchTripNotFound = &AppError{
		Code:    "SEARCH_TRIP_NOT_FOUND",
		Message: "Search trip not found in database",
	}

	ErrEventAlreadyProcessed = &AppError{
		Code:    "EVENT_ALREADY_PROCESSED",
		Message: "Event has already been processed (duplicate)",
	}
)

// IsNotFoundError checks if error is a not-found error
func IsNotFoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == "TRIP_NOT_FOUND" || appErr.Code == "USER_NOT_FOUND" || appErr.Code == "SEARCH_TRIP_NOT_FOUND"
	}
	return false
}

// IsServiceUnavailableError checks if error is a service unavailable error
func IsServiceUnavailableError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == "SERVICE_UNAVAILABLE"
	}
	return false
}

// IsTimeoutError checks if error is a timeout error
func IsTimeoutError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == "REQUEST_TIMEOUT"
	}
	return false
}

// WrapError wraps an error with context information
func WrapError(err error, context string) error {
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:    appErr.Code,
			Message: fmt.Sprintf("%s: %s", context, appErr.Message),
			Details: appErr.Details,
		}
	}
	return fmt.Errorf("%s: %w", context, err)
}
