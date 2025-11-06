package domain

// AppError representa un error de aplicación con código y mensaje
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	return e.Message
}

// Errores predefinidos de la aplicación
var (
	ErrTripNotFound         = &AppError{Code: "TRIP_NOT_FOUND", Message: "Trip not found"}
	ErrDriverNotFound       = &AppError{Code: "DRIVER_NOT_FOUND", Message: "Driver not found"}
	ErrNoSeatsAvailable     = &AppError{Code: "NO_SEATS_AVAILABLE", Message: "No seats available"}
	ErrOptimisticLockFailed = &AppError{Code: "OPTIMISTIC_LOCK_FAILED", Message: "Version conflict"}
	ErrUnauthorized         = &AppError{Code: "UNAUTHORIZED", Message: "Not authorized"}
	ErrPastDeparture        = &AppError{Code: "PAST_DEPARTURE", Message: "Departure must be in future"}
	ErrHasReservations      = &AppError{Code: "HAS_RESERVATIONS", Message: "Cannot modify trip with reservations"}
)
