package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Trip represents the external Trip DTO from trips-api
// This is the structure returned by trips-api GET /trips/:id endpoint
type Trip struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Driver and ownership
	DriverID int64 `json:"driver_id" bson:"driver_id"`

	// Location information (from trips-api with simple lat/lng format)
	Origin      TripLocation `json:"origin" bson:"origin"`
	Destination TripLocation `json:"destination" bson:"destination"`

	// Trip timing
	DepartureDatetime        time.Time `json:"departure_datetime" bson:"departure_datetime"`
	EstimatedArrivalDatetime time.Time `json:"estimated_arrival_datetime" bson:"estimated_arrival_datetime"`

	// Pricing and availability
	PricePerSeat   float64 `json:"price_per_seat" bson:"price_per_seat"`
	TotalSeats     int     `json:"total_seats" bson:"total_seats"`
	AvailableSeats int     `json:"available_seats" bson:"available_seats"`
	ReservedSeats  int     `json:"reserved_seats" bson:"reserved_seats"`

	// Optimistic locking for concurrency control
	AvailabilityVersion int `json:"availability_version" bson:"availability_version"`

	// Vehicle and preferences
	Car         Car         `json:"car" bson:"car"`
	Preferences Preferences `json:"preferences" bson:"preferences"`

	// Trip details
	Status      string `json:"status" bson:"status"` // draft, published, full, in_progress, completed, cancelled
	Description string `json:"description" bson:"description"`

	// Cancellation info (optional)
	CancelledAt        *time.Time `json:"cancelled_at,omitempty" bson:"cancelled_at,omitempty"`
	CancelledBy        *int64     `json:"cancelled_by,omitempty" bson:"cancelled_by,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// TripLocation represents location data from trips-api (with simple lat/lng coordinates)
type TripLocation struct {
	City        string            `json:"city"`
	Province    string            `json:"province"`
	Address     string            `json:"address"`
	Coordinates SimpleCoordinates `json:"coordinates"`
}

// SimpleCoordinates represents simple lat/lng coordinates from trips-api
type SimpleCoordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ToLocation converts TripLocation (from trips-api) to Location (for search-api storage)
func (tl *TripLocation) ToLocation() Location {
	return Location{
		City:        tl.City,
		Province:    tl.Province,
		Address:     tl.Address,
		Coordinates: NewGeoJSONPoint(tl.Coordinates.Lat, tl.Coordinates.Lng),
	}
}

// ToSearchTrip converts Trip DTO to SearchTrip with driver information
// Driver info must be fetched separately from users-api
// This method converts simple lat/lng coordinates to GeoJSON format required by MongoDB
func (t *Trip) ToSearchTrip(driver Driver) *SearchTrip {
	return &SearchTrip{
		TripID:                   t.ID.Hex(),
		DriverID:                 t.DriverID,
		Driver:                   driver,
		Origin:                   t.Origin.ToLocation(),
		Destination:              t.Destination.ToLocation(),
		DepartureDatetime:        t.DepartureDatetime,
		EstimatedArrivalDatetime: t.EstimatedArrivalDatetime,
		PricePerSeat:             t.PricePerSeat,
		TotalSeats:               t.TotalSeats,
		AvailableSeats:           t.AvailableSeats,
		Car:                      t.Car,
		Preferences:              t.Preferences,
		Status:                   t.Status,
		Description:              t.Description,
		SearchText:               buildSearchText(t, driver),
		CreatedAt:                t.CreatedAt,
		UpdatedAt:                t.UpdatedAt,
	}
}

// buildSearchText creates a concatenated text field for full-text search fallback
func buildSearchText(trip *Trip, driver Driver) string {
	return trip.Origin.City + " " +
		trip.Origin.Province + " " +
		trip.Destination.City + " " +
		trip.Destination.Province + " " +
		driver.Name + " " +
		trip.Car.Brand + " " +
		trip.Car.Model + " " +
		trip.Description
}
