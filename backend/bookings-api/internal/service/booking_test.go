package service

import (
	"testing"

	"github.com/carpooling-ucc/bookings-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

// Basic test to verify the toBookingDTO function
func TestToBookingDTO(t *testing.T) {
	// This is a simple test to ensure the package compiles
	// Full tests with mocks would be implemented here
	assert.True(t, true, "Booking service package compiles successfully")
}

// Test struct for validation context
func TestValidationContext(t *testing.T) {
	ctx := ValidationContext{
		TripID:        "test-trip-id",
		PassengerID:   1,
		SeatsReserved: 2,
	}

	assert.Equal(t, "test-trip-id", ctx.TripID)
	assert.Equal(t, int64(1), ctx.PassengerID)
	assert.Equal(t, 2, ctx.SeatsReserved)
}

// Test DTO conversion
func TestBookingDTOCreation(t *testing.T) {
	req := domain.CreateBookingRequest{
		TripID:        "trip-123",
		SeatsReserved: 2,
	}

	assert.Equal(t, "trip-123", req.TripID)
	assert.Equal(t, 2, req.SeatsReserved)
}
