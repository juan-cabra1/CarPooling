package testutil

import (
	"testing"
	"time"

	"github.com/juan-cabra1/CarPooling/backend/search-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

// AssertSearchTripEquals verifies that two SearchTrip objects are equal
func AssertSearchTripEquals(t *testing.T, expected, actual *domain.SearchTrip) {
	assert.Equal(t, expected.TripID, actual.TripID, "TripID should match")
	assert.Equal(t, expected.DriverID, actual.DriverID, "DriverID should match")
	assert.Equal(t, expected.Origin.City, actual.Origin.City, "Origin city should match")
	assert.Equal(t, expected.Destination.City, actual.Destination.City, "Destination city should match")
	assert.Equal(t, expected.AvailableSeats, actual.AvailableSeats, "AvailableSeats should match")
	assert.Equal(t, expected.PricePerSeat, actual.PricePerSeat, "PricePerSeat should match")
	assert.Equal(t, expected.Status, actual.Status, "Status should match")
}

// AssertLocationNear verifies that two locations are within a certain distance
func AssertLocationNear(t *testing.T, expected, actual domain.Location, deltaKm float64) {
	latDiff := abs(expected.Coordinates.Lat - actual.Coordinates.Lat)
	lngDiff := abs(expected.Coordinates.Lng - actual.Coordinates.Lng)

	// Rough approximation: 1 degree ≈ 111 km
	maxDelta := deltaKm / 111.0

	assert.LessOrEqual(t, latDiff, maxDelta, "Latitude should be within delta")
	assert.LessOrEqual(t, lngDiff, maxDelta, "Longitude should be within delta")
}

// AssertTimeNear verifies that two times are within a certain duration
func AssertTimeNear(t *testing.T, expected, actual time.Time, delta time.Duration) {
	diff := actual.Sub(expected)
	if diff < 0 {
		diff = -diff
	}
	assert.LessOrEqual(t, diff, delta, "Times should be within delta")
}

// AssertPopularityScoreValid verifies that popularity score is within valid range
func AssertPopularityScoreValid(t *testing.T, score float64) {
	assert.GreaterOrEqual(t, score, 0.0, "Popularity score should be >= 0")
	assert.LessOrEqual(t, score, 100.0, "Popularity score should be <= 100")
}

// AssertSearchTextContains verifies that search text contains expected terms
func AssertSearchTextContains(t *testing.T, searchText string, terms ...string) {
	for _, term := range terms {
		assert.Contains(t, searchText, term, "Search text should contain: %s", term)
	}
}

// AssertCacheKeyFormat verifies cache key format
func AssertCacheKeyFormat(t *testing.T, key, expectedPrefix string) {
	assert.Contains(t, key, expectedPrefix, "Cache key should contain prefix: %s", expectedPrefix)
}

// AssertDriverRatingValid verifies driver rating is in valid range
func AssertDriverRatingValid(t *testing.T, rating float64) {
	assert.GreaterOrEqual(t, rating, 0.0, "Driver rating should be >= 0")
	assert.LessOrEqual(t, rating, 5.0, "Driver rating should be <= 5")
}

// AssertTripAvailabilityValid verifies trip availability logic
func AssertTripAvailabilityValid(t *testing.T, trip *domain.SearchTrip) {
	assert.GreaterOrEqual(t, trip.AvailableSeats, 0, "Available seats should be >= 0")
	assert.GreaterOrEqual(t, trip.ReservedSeats, 0, "Reserved seats should be >= 0")
	totalSeats := trip.AvailableSeats + trip.ReservedSeats
	assert.Greater(t, totalSeats, 0, "Total seats should be > 0")
}

// AssertCoordinatesValid verifies coordinates are within valid ranges
func AssertCoordinatesValid(t *testing.T, coords domain.Coordinates) {
	assert.GreaterOrEqual(t, coords.Lat, -90.0, "Latitude should be >= -90")
	assert.LessOrEqual(t, coords.Lat, 90.0, "Latitude should be <= 90")
	assert.GreaterOrEqual(t, coords.Lng, -180.0, "Longitude should be >= -180")
	assert.LessOrEqual(t, coords.Lng, 180.0, "Longitude should be <= 180")
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
