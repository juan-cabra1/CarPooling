package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"search-api/internal/domain"

	"github.com/rs/zerolog/log"
)

// TripsClient defines the interface for communicating with trips-api
type TripsClient interface {
	GetTrip(ctx context.Context, tripID string) (*domain.Trip, error)
}

// tripsHTTPClient implements TripsClient using HTTP
type tripsHTTPClient struct {
	baseURL        string
	client         *http.Client
	maxRetries     int
	retryWaitTime  time.Duration
	circuitBreaker *CircuitBreaker
}

// NewTripsClient creates a new TripsClient with the given configuration
func NewTripsClient(config HTTPClientConfig) TripsClient {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryWaitTime == 0 {
		config.RetryWaitTime = 1 * time.Second
	}
	if config.CircuitBreaker == nil {
		config.CircuitBreaker = NewCircuitBreaker(5, 30*time.Second)
	}

	return &tripsHTTPClient{
		baseURL:        config.BaseURL,
		client:         CreateHTTPClient(config.Timeout),
		maxRetries:     config.MaxRetries,
		retryWaitTime:  config.RetryWaitTime,
		circuitBreaker: config.CircuitBreaker,
	}
}

// GetTrip fetches full trip details from trips-api
// Endpoint: GET /trips/:id
// Returns: Trip DTO with all details
func (c *tripsHTTPClient) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	url := fmt.Sprintf("%s/trips/%s", c.baseURL, tripID)

	log.Debug().
		Str("trip_id", tripID).
		Str("url", url).
		Msg("Fetching trip from trips-api")

	// Execute with circuit breaker
	var trip *domain.Trip
	err := c.circuitBreaker.Call(func() error {
		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return domain.WrapError(domain.ErrInvalidResponse, "failed to create HTTP request")
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "search-api/1.0")

		// Execute request with retry logic
		resp, err := DoRequestWithRetry(ctx, c.client, req, c.maxRetries, c.retryWaitTime)
		if err != nil {
			return err
		}

		// Handle 404 specifically for trips
		if resp.StatusCode == http.StatusNotFound {
			return domain.ErrTripNotFound
		}

		// Parse standard response
		var tripData domain.Trip
		if err := ParseStandardResponse(resp, &tripData); err != nil {
			return err
		}

		trip = &tripData
		return nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to fetch trip from trips-api")
		return nil, err
	}

	log.Info().
		Str("trip_id", tripID).
		Str("status", trip.Status).
		Int("available_seats", trip.AvailableSeats).
		Msg("Successfully fetched trip from trips-api")

	return trip, nil
}
