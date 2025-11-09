package domain

// Trip represents trip information from trips-api
// Used for HTTP responses when validating bookings
type Trip struct {
	ID             string  `json:"id"`
	DriverID       int64   `json:"driver_id"`
	AvailableSeats int     `json:"available_seats"`
	PricePerSeat   float64 `json:"price_per_seat"`
	Status         string  `json:"status"`
}

// Trip status constants
const (
	TripStatusDraft     = "draft"
	TripStatusPublished = "published"
	TripStatusCompleted = "completed"
	TripStatusCancelled = "cancelled"
)

// IsPublished checks if the trip is in published status
func (t *Trip) IsPublished() bool {
	return t.Status == TripStatusPublished
}

// HasAvailableSeats checks if the trip has enough available seats
func (t *Trip) HasAvailableSeats(requested int) bool {
	return t.AvailableSeats >= requested
}

// CalculateTotalPrice calculates the total price for requested seats
func (t *Trip) CalculateTotalPrice(seats int) float64 {
	return t.PricePerSeat * float64(seats)
}
