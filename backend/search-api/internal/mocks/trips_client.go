package mocks

import (
	"context"

	"github.com/juan-cabra1/CarPooling/backend/search-api/internal/domain"
)

// MockTripsClient is a mock implementation of the TripsClient interface
type MockTripsClient struct {
	GetTripFunc func(ctx context.Context, tripID string) (*domain.Trip, error)
}

// GetTrip calls the mocked GetTripFunc
func (m *MockTripsClient) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	if m.GetTripFunc != nil {
		return m.GetTripFunc(ctx, tripID)
	}
	return nil, nil
}
