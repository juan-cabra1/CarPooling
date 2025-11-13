package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"search-api/internal/domain"
	"search-api/internal/mocks"
	"search-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSearchTrips_CacheHit(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockPopularRouteRepo := &mocks.MockPopularRouteRepository{}
	mockSolr := &mocks.MockSolrClient{}
	mockTripsClient := &mocks.MockTripsClient{}
	mockUsersClient := &mocks.MockUsersClient{}

	service := NewSearchService(
		mockTripRepo,
		mockPopularRouteRepo,
		mockCache,
		mockSolr,
		mockTripsClient,
		mockUsersClient,
	)

	query := testutil.CreateTestSearchQuery()
	expectedResponse := &domain.SearchResponse{
		Trips:      []*domain.SearchTrip{testutil.CreateTestSearchTrip("trip-1")},
		Total:      1,
		Page:       1,
		Limit:      20,
		TotalPages: 1,
	}

	// Mock cache hit
	responseJSON, _ := json.Marshal(expectedResponse)
	mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
		assert.Contains(t, key, "search:query:")
		return string(responseJSON), nil
	}

	// Mock popular route tracking (fire-and-forget)
	mockPopularRouteRepo.IncrementSearchCountFunc = func(ctx context.Context, originCity, destinationCity string) error {
		return nil
	}

	// Execute
	result, err := service.SearchTrips(context.Background(), query)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Total, result.Total)
	assert.Len(t, result.Trips, 1)
}

func TestSearchTrips_CacheMiss_MongoDBFallback(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{}
	mockTripRepo := &mocks.MockTripRepository{}
	mockPopularRouteRepo := &mocks.MockPopularRouteRepository{}
	mockTripsClient := &mocks.MockTripsClient{}
	mockUsersClient := &mocks.MockUsersClient{}

	service := NewSearchService(
		mockTripRepo,
		mockPopularRouteRepo,
		mockCache,
		nil, // No Solr client
		mockTripsClient,
		mockUsersClient,
	)

	query := testutil.CreateTestSearchQuery()
	trips := []*domain.SearchTrip{
		testutil.CreateTestSearchTrip("trip-1"),
		testutil.CreateTestSearchTrip("trip-2"),
	}

	// Mock cache miss
	mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
		return "", errors.New("cache miss")
	}

	// Mock MongoDB search
	mockTripRepo.SearchFunc = func(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.SearchTrip, int64, error) {
		assert.Equal(t, "published", filters["status"])
		assert.Equal(t, query.OriginCity, filters["origin.city"])
		assert.Equal(t, query.DestinationCity, filters["destination.city"])
		return trips, 2, nil
	}

	// Mock cache set
	mockCache.SetFunc = func(ctx context.Context, key string, value string, ttl time.Duration) error {
		assert.Contains(t, key, "search:query:")
		assert.Equal(t, 15*time.Minute, ttl)
		return nil
	}

	// Mock popular route tracking
	mockPopularRouteRepo.IncrementSearchCountFunc = func(ctx context.Context, originCity, destinationCity string) error {
		assert.Equal(t, query.OriginCity, originCity)
		assert.Equal(t, query.DestinationCity, destinationCity)
		return nil
	}

	// Execute
	result, err := service.SearchTrips(context.Background(), query)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Trips, 2)
	assert.Equal(t, 1, result.TotalPages)
}

func TestSearchTrips_ValidationError(t *testing.T) {
	// Setup
	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Create invalid query (negative page)
	query := &domain.SearchQuery{
		Page:  -1,
		Limit: 20,
	}

	// Execute
	result, err := service.SearchTrips(context.Background(), query)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid query")
}

func TestSearchTrips_FilterCombinations(t *testing.T) {
	tests := []struct {
		name          string
		query         *domain.SearchQuery
		expectedFilters map[string]interface{}
	}{
		{
			name: "City filters",
			query: &domain.SearchQuery{
				OriginCity:      "Bogotá",
				DestinationCity: "Medellín",
				Page:            1,
				Limit:           20,
			},
			expectedFilters: map[string]interface{}{
				"status":           "published",
				"origin.city":      "Bogotá",
				"destination.city": "Medellín",
			},
		},
		{
			name: "Price and seats filters",
			query: func() *domain.SearchQuery {
				maxPrice := 50000.0
				minSeats := 2
				return &domain.SearchQuery{
					MaxPrice: &maxPrice,
					MinSeats: &minSeats,
					Page:     1,
					Limit:    20,
				}
			}(),
			expectedFilters: map[string]interface{}{
				"status":           "published",
				"price_per_seat":   map[string]interface{}{"$lte": 50000.0},
				"available_seats":  map[string]interface{}{"$gte": 2},
			},
		},
		{
			name: "Preference filters",
			query: func() *domain.SearchQuery {
				petsAllowed := true
				smokingAllowed := false
				return &domain.SearchQuery{
					PetsAllowed:    &petsAllowed,
					SmokingAllowed: &smokingAllowed,
					Page:           1,
					Limit:          20,
				}
			}(),
			expectedFilters: map[string]interface{}{
				"status":                       "published",
				"preferences.pets_allowed":     true,
				"preferences.smoking_allowed":  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := &mocks.MockCache{
				GetFunc: func(ctx context.Context, key string) (string, error) {
					return "", errors.New("cache miss")
				},
				SetFunc: func(ctx context.Context, key string, value string, ttl time.Duration) error {
					return nil
				},
			}

			mockTripRepo := &mocks.MockTripRepository{}
			mockTripRepo.SearchFunc = func(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.SearchTrip, int64, error) {
				// Verify expected filters are present
				for key, expectedValue := range tt.expectedFilters {
					actualValue, exists := filters[key]
					assert.True(t, exists, "Filter %s should exist", key)
					assert.Equal(t, expectedValue, actualValue, "Filter %s should match", key)
				}
				return []*domain.SearchTrip{}, 0, nil
			}

			service := NewSearchService(
				mockTripRepo,
				&mocks.MockPopularRouteRepository{
					IncrementSearchCountFunc: func(ctx context.Context, originCity, destinationCity string) error {
						return nil
					},
				},
				mockCache,
				nil,
				&mocks.MockTripsClient{},
				&mocks.MockUsersClient{},
			)

			_, err := service.SearchTrips(context.Background(), tt.query)
			assert.NoError(t, err)
		})
	}
}

func TestSearchByLocation_ValidCoordinates(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("cache miss")
		},
		SetFunc: func(ctx context.Context, key string, value string, ttl time.Duration) error {
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{}
	trips := []domain.SearchTrip{
		*testutil.CreateTestSearchTrip("trip-1"),
		*testutil.CreateTestSearchTrip("trip-2"),
	}

	mockTripRepo.SearchByLocationFunc = func(ctx context.Context, lat, lng, radiusKm float64, filters map[string]interface{}) ([]domain.SearchTrip, error) {
		assert.Equal(t, 4.7110, lat)
		assert.Equal(t, -74.0721, lng)
		assert.Equal(t, 10.0, radiusKm)
		return trips, nil
	}

	service := NewSearchService(
		mockTripRepo,
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.SearchByLocation(context.Background(), 4.7110, -74.0721, 10, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Trips, 2)
}

func TestSearchByLocation_InvalidCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		lat       float64
		lng       float64
		radiusKm  int
		expectErr string
	}{
		{
			name:      "Invalid latitude (too high)",
			lat:       91.0,
			lng:       -74.0721,
			radiusKm:  10,
			expectErr: "latitude must be between -90 and 90",
		},
		{
			name:      "Invalid latitude (too low)",
			lat:       -91.0,
			lng:       -74.0721,
			radiusKm:  10,
			expectErr: "latitude must be between -90 and 90",
		},
		{
			name:      "Invalid longitude (too high)",
			lat:       4.7110,
			lng:       181.0,
			radiusKm:  10,
			expectErr: "longitude must be between -180 and 180",
		},
		{
			name:      "Invalid radius (negative)",
			lat:       4.7110,
			lng:       -74.0721,
			radiusKm:  -5,
			expectErr: "radius must be positive",
		},
		{
			name:      "Invalid radius (zero)",
			lat:       4.7110,
			lng:       -74.0721,
			radiusKm:  0,
			expectErr: "radius must be positive",
		},
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.SearchByLocation(context.Background(), tt.lat, tt.lng, tt.radiusKm, nil)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.expectErr)
		})
	}
}

func TestSearchByLocation_CacheHit(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{}
	expectedResponse := &domain.SearchResponse{
		Trips: []domain.SearchTrip{
			*testutil.CreateTestSearchTrip("trip-1"),
		},
		Total:      1,
		Page:       1,
		Limit:      1,
		TotalPages: 1,
	}

	responseJSON, _ := json.Marshal(expectedResponse)
	mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
		assert.Contains(t, key, "search:location:")
		return string(responseJSON), nil
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.SearchByLocation(context.Background(), 4.7110, -74.0721, 10, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)
}

func TestGetTrip_CacheHit(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{}
	expectedTrip := testutil.CreateTestSearchTrip("trip-123")
	tripJSON, _ := json.Marshal(expectedTrip)

	mockCache.GetFunc = func(ctx context.Context, key string) (string, error) {
		assert.Equal(t, "trip:trip-123", key)
		return string(tripJSON), nil
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.GetTrip(context.Background(), "trip-123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTrip.TripID, result.TripID)
}

func TestGetTrip_CacheMiss_FetchFromMongoDB(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("cache miss")
		},
		SetFunc: func(ctx context.Context, key string, value string, ttl time.Duration) error {
			assert.Contains(t, key, "trip:")
			return nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{}
	expectedTrip := testutil.CreateTestSearchTrip("trip-123")

	mockTripRepo.FindByIDFunc = func(ctx context.Context, id primitive.ObjectID) (*domain.SearchTrip, error) {
		return expectedTrip, nil
	}

	service := NewSearchService(
		mockTripRepo,
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.GetTrip(context.Background(), "trip-123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTrip.TripID, result.TripID)
}

func TestGetTrip_NotFound(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("cache miss")
		},
	}

	mockTripRepo := &mocks.MockTripRepository{}
	mockTripRepo.FindByIDFunc = func(ctx context.Context, id primitive.ObjectID) (*domain.SearchTrip, error) {
		return nil, nil
	}

	service := NewSearchService(
		mockTripRepo,
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.GetTrip(context.Background(), "nonexistent-trip")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "trip not found")
}

func TestGetPopularRoutes(t *testing.T) {
	// Setup
	mockPopularRouteRepo := &mocks.MockPopularRouteRepository{}
	expectedRoutes := []domain.PopularRoute{
		testutil.CreateTestPopularRoute("Bogotá", "Medellín", 150),
		testutil.CreateTestPopularRoute("Medellín", "Cali", 120),
		testutil.CreateTestPopularRoute("Cali", "Barranquilla", 90),
	}

	mockPopularRouteRepo.GetTopRoutesFunc = func(ctx context.Context, limit int) ([]domain.PopularRoute, error) {
		assert.Equal(t, 10, limit)
		return expectedRoutes, nil
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		mockPopularRouteRepo,
		&mocks.MockCache{},
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	result, err := service.GetPopularRoutes(context.Background(), 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Bogotá", result[0].OriginCity)
	assert.Equal(t, "Medellín", result[0].DestinationCity)
	assert.Equal(t, 150, result[0].SearchCount)
}

func TestGetAutocomplete(t *testing.T) {
	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute - currently returns empty array
	result, err := service.GetAutocomplete(context.Background(), "Bog", 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result) // TODO: Will be populated when implemented
}

func TestInvalidateCache(t *testing.T) {
	// Setup
	mockCache := &mocks.MockCache{}
	deleteCalled := false

	mockCache.DeleteFunc = func(ctx context.Context, key string) error {
		assert.Equal(t, "trip:trip-123", key)
		deleteCalled = true
		return nil
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		mockCache,
		nil,
		&mocks.MockTripsClient{},
		&mocks.MockUsersClient{},
	)

	// Execute
	err := service.InvalidateCache(context.Background(), "trip-123")

	// Assert
	require.NoError(t, err)
	assert.True(t, deleteCalled, "Cache delete should have been called")
}

func TestBuildSearchText(t *testing.T) {
	trip := testutil.CreateTestSearchTrip("trip-123")
	trip.Origin.City = "Bogotá"
	trip.Origin.Province = "Cundinamarca"
	trip.Destination.City = "Medellín"
	trip.Destination.Province = "Antioquia"
	trip.Description = "Comfortable trip"
	trip.Driver.Name = "John Doe"

	searchText := BuildSearchText(trip)

	// Assert all components are present
	testutil.AssertSearchTextContains(t, searchText, "Bogotá", "Cundinamarca", "Medellín", "Antioquia", "Comfortable trip", "John Doe")
}

func TestCalculatePopularityScore(t *testing.T) {
	tests := []struct {
		name          string
		trip          *domain.SearchTrip
		expectedMin   float64
		expectedMax   float64
	}{
		{
			name: "High occupancy, high rating, experienced driver",
			trip: func() *domain.SearchTrip {
				trip := testutil.CreateTestSearchTrip("trip-1")
				trip.TotalSeats = 4
				trip.AvailableSeats = 0 // Fully booked
				trip.Driver.Rating = 5.0
				trip.Driver.TotalTrips = 100
				trip.CreatedAt = time.Now()
				return trip
			}(),
			expectedMin: 80.0,
			expectedMax: 100.0,
		},
		{
			name: "Low occupancy, low rating, new driver",
			trip: func() *domain.SearchTrip {
				trip := testutil.CreateTestSearchTrip("trip-2")
				trip.TotalSeats = 4
				trip.AvailableSeats = 4 // Empty
				trip.Driver.Rating = 1.0
				trip.Driver.TotalTrips = 1
				trip.CreatedAt = time.Now()
				return trip
			}(),
			expectedMin: 5.0,
			expectedMax: 20.0,
		},
		{
			name: "Medium occupancy, medium rating",
			trip: func() *domain.SearchTrip {
				trip := testutil.CreateTestSearchTrip("trip-3")
				trip.TotalSeats = 4
				trip.AvailableSeats = 2 // 50% occupied
				trip.Driver.Rating = 3.5
				trip.Driver.TotalTrips = 50
				trip.CreatedAt = time.Now()
				return trip
			}(),
			expectedMin: 40.0,
			expectedMax: 60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculatePopularityScore(tt.trip)
			testutil.AssertPopularityScoreValid(t, score)
			assert.GreaterOrEqual(t, score, tt.expectedMin, "Score should be >= expected minimum")
			assert.LessOrEqual(t, score, tt.expectedMax, "Score should be <= expected maximum")
		})
	}
}

func TestDenormalizeTrip_Success(t *testing.T) {
	// Setup
	tripID := primitive.NewObjectID().Hex()
	testTrip := testutil.CreateTestTrip(tripID)
	testUser := testutil.CreateTestUser(123)

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			assert.Equal(t, tripID, id)
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, userID string) (*domain.User, error) {
			assert.Equal(t, "123", userID)
			return testUser, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			assert.Equal(t, tripID, trip.TripID)
			assert.NotEmpty(t, trip.SearchText)
			assert.Greater(t, trip.PopularityScore, 0.0)
			return nil
		},
	}

	mockSolr := &mocks.MockSolrClient{
		IndexFunc: func(trip *domain.SearchTrip) error {
			return nil
		},
	}

	service := NewSearchService(
		mockTripRepo,
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		mockSolr,
		mockTripsClient,
		mockUsersClient,
	)

	// Execute
	err := service.DenormalizeTrip(context.Background(), tripID)

	// Assert
	require.NoError(t, err)
}

func TestDenormalizeTrip_TripsAPIFailure(t *testing.T) {
	// Setup
	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return nil, errors.New("trips-api unavailable")
		},
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		nil,
		mockTripsClient,
		&mocks.MockUsersClient{},
	)

	// Execute
	err := service.DenormalizeTrip(context.Background(), "trip-123")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch trip")
}

func TestDenormalizeTrip_UsersAPIFailure(t *testing.T) {
	// Setup
	tripID := primitive.NewObjectID().Hex()
	testTrip := testutil.CreateTestTrip(tripID)

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, userID string) (*domain.User, error) {
			return nil, errors.New("users-api unavailable")
		},
	}

	service := NewSearchService(
		&mocks.MockTripRepository{},
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		nil,
		mockTripsClient,
		mockUsersClient,
	)

	// Execute
	err := service.DenormalizeTrip(context.Background(), tripID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch driver")
}

func TestDenormalizeTrip_SolrFailureDoesNotBlock(t *testing.T) {
	// Setup - Solr failure should not prevent the operation from succeeding
	tripID := primitive.NewObjectID().Hex()
	testTrip := testutil.CreateTestTrip(tripID)
	testUser := testutil.CreateTestUser(123)

	mockTripsClient := &mocks.MockTripsClient{
		GetTripFunc: func(ctx context.Context, id string) (*domain.Trip, error) {
			return testTrip, nil
		},
	}

	mockUsersClient := &mocks.MockUsersClient{
		GetUserFunc: func(ctx context.Context, userID string) (*domain.User, error) {
			return testUser, nil
		},
	}

	mockTripRepo := &mocks.MockTripRepository{
		CreateFunc: func(ctx context.Context, trip *domain.SearchTrip) error {
			return nil
		},
	}

	mockSolr := &mocks.MockSolrClient{
		IndexFunc: func(trip *domain.SearchTrip) error {
			return errors.New("solr unavailable")
		},
	}

	service := NewSearchService(
		mockTripRepo,
		&mocks.MockPopularRouteRepository{},
		&mocks.MockCache{},
		mockSolr,
		mockTripsClient,
		mockUsersClient,
	)

	// Execute
	err := service.DenormalizeTrip(context.Background(), tripID)

	// Assert - Should succeed despite Solr failure
	require.NoError(t, err)
}
