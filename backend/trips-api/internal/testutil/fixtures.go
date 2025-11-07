package testutil

import (
	"time"

	"trips-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewTestTrip creates a valid Trip for testing purposes
func NewTestTrip(driverID int64) *domain.Trip {
	now := time.Now()
	departure := now.Add(24 * time.Hour) // Tomorrow
	arrival := departure.Add(2 * time.Hour)

	return &domain.Trip{
		ID:                       primitive.NewObjectID(),
		DriverID:                 driverID,
		Origin:                   NewTestLocation("Buenos Aires"),
		Destination:              NewTestLocation("Rosario"),
		DepartureDatetime:        departure,
		EstimatedArrivalDatetime: arrival,
		PricePerSeat:             1500.0,
		TotalSeats:               4,
		ReservedSeats:            0,
		AvailableSeats:           4,
		AvailabilityVersion:      1,
		Car:                      NewTestCar(),
		Preferences:              NewTestPreferences(),
		Status:                   "published",
		Description:              "Test trip description",
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

// NewTestLocation creates a valid Location for testing
func NewTestLocation(city string) domain.Location {
	return domain.Location{
		City:     city,
		Province: "Buenos Aires",
		Address:  "Av. Test 123",
		Coordinates: domain.Coordinates{
			Lat: -34.6037,
			Lng: -58.3816,
		},
	}
}

// NewTestCar creates a valid Car for testing
func NewTestCar() domain.Car {
	return domain.Car{
		Brand: "Toyota",
		Model: "Corolla",
		Year:  2020,
		Color: "Blanco",
		Plate: "ABC123",
	}
}

// NewTestPreferences creates default Preferences for testing
func NewTestPreferences() domain.Preferences {
	return domain.Preferences{
		PetsAllowed:    true,
		SmokingAllowed: false,
		MusicAllowed:   true,
	}
}

// NewTestCreateTripRequest creates a valid CreateTripRequest for testing
func NewTestCreateTripRequest() domain.CreateTripRequest {
	now := time.Now()
	departure := now.Add(24 * time.Hour)
	arrival := departure.Add(2 * time.Hour)

	return domain.CreateTripRequest{
		Origin:                   NewTestLocation("Buenos Aires"),
		Destination:              NewTestLocation("Rosario"),
		DepartureDatetime:        departure.Format(time.RFC3339),
		EstimatedArrivalDatetime: arrival.Format(time.RFC3339),
		PricePerSeat:             1500.0,
		TotalSeats:               4,
		Car:                      NewTestCar(),
		Preferences:              NewTestPreferences(),
		Description:              "Test trip",
	}
}

// NewTestProcessedEvent creates a valid ProcessedEvent for testing
func NewTestProcessedEvent(eventID, eventType string) *domain.ProcessedEvent {
	return &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   eventType,
		ProcessedAt: time.Now(),
	}
}

// WithDriverID sets the driver ID for a trip
func WithDriverID(trip *domain.Trip, driverID int64) *domain.Trip {
	trip.DriverID = driverID
	return trip
}

// WithAvailableSeats sets available seats for a trip
func WithAvailableSeats(trip *domain.Trip, seats int) *domain.Trip {
	trip.AvailableSeats = seats
	trip.ReservedSeats = trip.TotalSeats - seats
	return trip
}

// WithStatus sets the status for a trip
func WithStatus(trip *domain.Trip, status string) *domain.Trip {
	trip.Status = status
	return trip
}

// WithVersion sets the availability version for a trip
func WithVersion(trip *domain.Trip, version int) *domain.Trip {
	trip.AvailabilityVersion = version
	return trip
}
