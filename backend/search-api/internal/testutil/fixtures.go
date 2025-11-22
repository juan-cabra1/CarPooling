package testutil

import (
	"time"

	"github.com/juan-cabra1/CarPooling/backend/search-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTestTrip creates a test Trip domain object with default values
func CreateTestTrip(tripID string) *domain.Trip {
	departureTime := time.Now().Add(24 * time.Hour)
	objID, _ := primitive.ObjectIDFromHex(tripID)
	return &domain.Trip{
		ID:                       objID,
		DriverID:                 123,
		Origin:                   CreateTestLocation("Bogotá", 4.7110, -74.0721),
		Destination:              CreateTestLocation("Medellín", 6.2442, -75.5812),
		DepartureDatetime:        departureTime,
		EstimatedArrivalDatetime: departureTime.Add(6 * time.Hour),
		AvailableSeats:           3,
		ReservedSeats:            1,
		TotalSeats:               4,
		PricePerSeat:             50000,
		Car:                      CreateTestCar(),
		Preferences:              CreateTestPreferences(),
		Status:                   "published",
		Description:              "Trip to Medellín for business meeting",
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
}

// CreateTestSearchTrip creates a test SearchTrip domain object
func CreateTestSearchTrip(tripID string) *domain.SearchTrip {
	departureTime := time.Now().Add(24 * time.Hour)
	trip := &domain.SearchTrip{
		ID:                       primitive.NewObjectID(),
		TripID:                   tripID,
		DriverID:                 123,
		Origin:                   CreateTestLocation("Bogotá", 4.7110, -74.0721),
		Destination:              CreateTestLocation("Medellín", 6.2442, -75.5812),
		DepartureDatetime:        departureTime,
		EstimatedArrivalDatetime: departureTime.Add(6 * time.Hour),
		AvailableSeats:           3,
		TotalSeats:               4,
		PricePerSeat:             50000,
		Car:                      CreateTestCar(),
		Preferences:              CreateTestPreferences(),
		Driver:                   CreateTestDriver(),
		Status:                   "published",
		Description:              "Trip to Medellín for business meeting",
		SearchText:               "Bogotá Medellín Trip to Medellín for business meeting",
		PopularityScore:          75.5,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
	return trip
}

// CreateTestLocation creates a test Location with NEW GeoJSON format
func CreateTestLocation(city string, lat, lng float64) domain.Location {
	return domain.Location{
		City:        city,
		Province:    "", // Empty by default, can be set explicitly if needed
		Address:     "Main Street 123, " + city,
		Coordinates: domain.NewGeoJSONPoint(lat, lng),
	}
}

// CreateTestCar creates a test Car
func CreateTestCar() domain.Car {
	return domain.Car{
		Brand:        "Toyota",
		Model:        "Corolla",
		Year:         2020,
		Color:        "White",
		LicensePlate: "ABC123",
	}
}

// CreateTestPreferences creates test Preferences
func CreateTestPreferences() domain.Preferences {
	petsAllowed := false
	smokingAllowed := false
	musicAllowed := true
	return domain.Preferences{
		PetsAllowed:    &petsAllowed,
		SmokingAllowed: &smokingAllowed,
		MusicAllowed:   &musicAllowed,
	}
}

// CreateTestDriver creates a test Driver
func CreateTestDriver() domain.Driver {
	return domain.Driver{
		ID:         123,
		Name:       "John Doe",
		Email:      "john.doe@example.com",
		PhotoURL:   "https://example.com/photo.jpg",
		Rating:     4.7,
		TotalTrips: 150,
	}
}

// CreateTestUser creates a test User
func CreateTestUser(userID int64) *domain.User {
	return &domain.User{
		ID:                     userID,
		Name:                   "John Doe",
		Email:                  "john.doe@example.com",
		PhotoURL:               "https://example.com/photo.jpg",
		AverageRatingAsDriver:  4.7,
		TotalTripsAsDriver:     150,
		AverageRatingAsPassenger: 4.8,
		TotalTripsAsPassenger:  75,
		CreatedAt:              time.Now(),
	}
}

// CreateTestProcessedEvent creates a test ProcessedEvent
func CreateTestProcessedEvent(eventID, eventType, tripID string) *domain.ProcessedEvent {
	return &domain.ProcessedEvent{
		EventID:     eventID,
		EventType:   eventType,
		TripID:      tripID,
		ProcessedAt: time.Now(),
	}
}

// CreateTestPopularRoute creates a test PopularRoute
func CreateTestPopularRoute(originCity, destinationCity string, count int) domain.PopularRoute {
	return domain.PopularRoute{
		ID:              primitive.NewObjectID(),
		OriginCity:      originCity,
		DestinationCity: destinationCity,
		SearchCount:     count,
		LastSearchedAt:  time.Now(),
	}
}

// CreateTestSearchQuery creates a test SearchQuery with NEW structure
func CreateTestSearchQuery() *domain.SearchQuery {
	petsAllowed := false
	smokingAllowed := false
	musicAllowed := true
	departureDate := time.Now()

	return &domain.SearchQuery{
		Origin: &domain.Location{
			City:        "Bogotá",
			Province:    "Cundinamarca",
			Address:     "Test Address",
			Coordinates: domain.NewGeoJSONPoint(4.7110, -74.0721),
		},
		Destination: &domain.Location{
			City:        "Medellín",
			Province:    "Antioquia",
			Address:     "Test Address",
			Coordinates: domain.NewGeoJSONPoint(6.2442, -75.5812),
		},
		OriginRadius:      0, // No radius by default
		DestinationRadius: 0,
		DepartureDate:     &departureDate,
		MinSeats:          1,
		MaxPrice:          100000,
		MinDriverRating:   4.0,
		PetsAllowed:       &petsAllowed,
		SmokingAllowed:    &smokingAllowed,
		MusicAllowed:      &musicAllowed,
		Page:              1,
		Limit:             20,
		SortBy:            "popularity",
	}
}

// CreateMultipleTestTrips creates multiple test trips with different characteristics
func CreateMultipleTestTrips(count int) []domain.SearchTrip {
	trips := make([]domain.SearchTrip, count)
	cities := []struct{ name string; lat, lng float64 }{
		{"Bogotá", 4.7110, -74.0721},
		{"Medellín", 6.2442, -75.5812},
		{"Cali", 3.4516, -76.5320},
		{"Barranquilla", 10.9685, -74.7813},
		{"Cartagena", 10.3910, -75.4794},
	}

	for i := 0; i < count; i++ {
		originIdx := i % len(cities)
		destIdx := (i + 1) % len(cities)

		departureTime := time.Now().Add(time.Duration(i*24) * time.Hour)
		trip := CreateTestSearchTrip(primitive.NewObjectID().Hex())
		trip.TripID = primitive.NewObjectID().Hex()
		trip.Origin = CreateTestLocation(cities[originIdx].name, cities[originIdx].lat, cities[originIdx].lng)
		trip.Destination = CreateTestLocation(cities[destIdx].name, cities[destIdx].lat, cities[destIdx].lng)
		trip.DepartureDatetime = departureTime
		trip.AvailableSeats = (i % 3) + 1
		trip.TotalSeats = 4
		trip.PricePerSeat = float64(30000 + (i * 10000))
		trip.PopularityScore = float64(50 + (i * 5))

		trips[i] = *trip
	}

	return trips
}
