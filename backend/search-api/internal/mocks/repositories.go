package mocks

import (
	"context"

	"search-api/internal/domain"
)

// MockTripRepository is a mock implementation of TripRepository
type MockTripRepository struct {
	CreateFunc                     func(ctx context.Context, trip *domain.SearchTrip) error
	FindByIDFunc                   func(ctx context.Context, id string) (*domain.SearchTrip, error)
	FindByTripIDFunc               func(ctx context.Context, tripID string) (*domain.SearchTrip, error)
	UpdateFunc                     func(ctx context.Context, trip *domain.SearchTrip) error
	UpdateStatusFunc               func(ctx context.Context, id string, status string) error
	UpdateStatusByTripIDFunc       func(ctx context.Context, tripID string, status string) error
	UpdateAvailabilityFunc         func(ctx context.Context, id string, availableSeats int) error
	UpdateAvailabilityByTripIDFunc func(ctx context.Context, tripID string, availableSeats, reservedSeats int, status string) error
	DeleteByTripIDFunc             func(ctx context.Context, tripID string) error
	SearchFunc                     func(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.SearchTrip, int64, error)
	SearchByLocationFunc           func(ctx context.Context, lat, lng float64, radiusKm int, additionalFilters map[string]interface{}) ([]*domain.SearchTrip, error)
	SearchByRouteFunc              func(ctx context.Context, originCity, destinationCity string, filters map[string]interface{}) ([]*domain.SearchTrip, error)
}

// Create calls the mocked CreateFunc
func (m *MockTripRepository) Create(ctx context.Context, trip *domain.SearchTrip) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, trip)
	}
	return nil
}

// FindByID calls the mocked FindByIDFunc
func (m *MockTripRepository) FindByID(ctx context.Context, id string) (*domain.SearchTrip, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

// FindByTripID calls the mocked FindByTripIDFunc
func (m *MockTripRepository) FindByTripID(ctx context.Context, tripID string) (*domain.SearchTrip, error) {
	if m.FindByTripIDFunc != nil {
		return m.FindByTripIDFunc(ctx, tripID)
	}
	return nil, nil
}

// Update calls the mocked UpdateFunc
func (m *MockTripRepository) Update(ctx context.Context, trip *domain.SearchTrip) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, trip)
	}
	return nil
}

// UpdateStatus calls the mocked UpdateStatusFunc
func (m *MockTripRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

// UpdateStatusByTripID calls the mocked UpdateStatusByTripIDFunc
func (m *MockTripRepository) UpdateStatusByTripID(ctx context.Context, tripID string, status string) error {
	if m.UpdateStatusByTripIDFunc != nil {
		return m.UpdateStatusByTripIDFunc(ctx, tripID, status)
	}
	return nil
}

// UpdateAvailability calls the mocked UpdateAvailabilityFunc
func (m *MockTripRepository) UpdateAvailability(ctx context.Context, id string, availableSeats int) error {
	if m.UpdateAvailabilityFunc != nil {
		return m.UpdateAvailabilityFunc(ctx, id, availableSeats)
	}
	return nil
}

// UpdateAvailabilityByTripID calls the mocked UpdateAvailabilityByTripIDFunc
func (m *MockTripRepository) UpdateAvailabilityByTripID(ctx context.Context, tripID string, availableSeats, reservedSeats int, status string) error {
	if m.UpdateAvailabilityByTripIDFunc != nil {
		return m.UpdateAvailabilityByTripIDFunc(ctx, tripID, availableSeats, reservedSeats, status)
	}
	return nil
}

// DeleteByTripID calls the mocked DeleteByTripIDFunc
func (m *MockTripRepository) DeleteByTripID(ctx context.Context, tripID string) error {
	if m.DeleteByTripIDFunc != nil {
		return m.DeleteByTripIDFunc(ctx, tripID)
	}
	return nil
}

// Search calls the mocked SearchFunc
func (m *MockTripRepository) Search(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.SearchTrip, int64, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, filters, page, limit)
	}
	return []*domain.SearchTrip{}, 0, nil
}

// SearchByLocation calls the mocked SearchByLocationFunc
func (m *MockTripRepository) SearchByLocation(ctx context.Context, lat, lng float64, radiusKm int, additionalFilters map[string]interface{}) ([]*domain.SearchTrip, error) {
	if m.SearchByLocationFunc != nil {
		return m.SearchByLocationFunc(ctx, lat, lng, radiusKm, additionalFilters)
	}
	return []*domain.SearchTrip{}, nil
}

// SearchByRoute calls the mocked SearchByRouteFunc
func (m *MockTripRepository) SearchByRoute(ctx context.Context, originCity, destinationCity string, filters map[string]interface{}) ([]*domain.SearchTrip, error) {
	if m.SearchByRouteFunc != nil {
		return m.SearchByRouteFunc(ctx, originCity, destinationCity, filters)
	}
	return []*domain.SearchTrip{}, nil
}

// MockEventRepository is a mock implementation of EventRepository
type MockEventRepository struct {
	IsEventProcessedFunc  func(ctx context.Context, eventID string) (bool, error)
	MarkEventProcessedFunc func(ctx context.Context, event *domain.ProcessedEvent) error
}

// IsEventProcessed calls the mocked IsEventProcessedFunc
func (m *MockEventRepository) IsEventProcessed(ctx context.Context, eventID string) (bool, error) {
	if m.IsEventProcessedFunc != nil {
		return m.IsEventProcessedFunc(ctx, eventID)
	}
	return false, nil
}

// MarkEventProcessed calls the mocked MarkEventProcessedFunc
func (m *MockEventRepository) MarkEventProcessed(ctx context.Context, event *domain.ProcessedEvent) error {
	if m.MarkEventProcessedFunc != nil {
		return m.MarkEventProcessedFunc(ctx, event)
	}
	return nil
}

// MockPopularRouteRepository is a mock implementation of PopularRouteRepository
type MockPopularRouteRepository struct {
	GetTopRoutesFunc         func(ctx context.Context, limit int) ([]domain.PopularRoute, error)
	IncrementSearchCountFunc func(ctx context.Context, originCity, destinationCity string) error
}

// GetTopRoutes calls the mocked GetTopRoutesFunc
func (m *MockPopularRouteRepository) GetTopRoutes(ctx context.Context, limit int) ([]domain.PopularRoute, error) {
	if m.GetTopRoutesFunc != nil {
		return m.GetTopRoutesFunc(ctx, limit)
	}
	return []domain.PopularRoute{}, nil
}

// IncrementSearchCount calls the mocked IncrementSearchCountFunc
func (m *MockPopularRouteRepository) IncrementSearchCount(ctx context.Context, originCity, destinationCity string) error {
	if m.IncrementSearchCountFunc != nil {
		return m.IncrementSearchCountFunc(ctx, originCity, destinationCity)
	}
	return nil
}
