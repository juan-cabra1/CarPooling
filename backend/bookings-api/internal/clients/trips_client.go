package clients

import (
	"bookings-api/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// TripsClient defines the interface for interacting with trips-api
type TripsClient interface {
	GetTrip(ctx context.Context, tripID string) (*domain.Trip, error)
}

// tripsHTTPClient implements TripsClient using HTTP calls
type tripsHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTripsClient creates a new HTTP client for trips-api
func NewTripsClient(baseURL string) TripsClient {
	return &tripsHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // 5 second timeout for external calls
		},
	}
}

// tripsAPIResponse represents the standardized API response wrapper from trips-api
type tripsAPIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// GetTrip retrieves trip details from trips-api
func (c *tripsHTTPClient) GetTrip(ctx context.Context, tripID string) (*domain.Trip, error) {
	// 1. Build URL
	url := fmt.Sprintf("%s/trips/%s", c.baseURL, tripID)

	log.Debug().
		Str("url", url).
		Str("trip_id", tripID).
		Msg("Calling trips-api to get trip details")

	// 2. Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Set headers
	req.Header.Set("Content-Type", "application/json")

	// 4. Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to call trips-api")
		return nil, domain.ErrTripsAPIUnavailable.WithDetails(map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	// 5. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Str("trip_id", tripID).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug().
		Str("trip_id", tripID).
		Int("status_code", resp.StatusCode).
		Str("response_body", string(body)).
		Msg("Received response from trips-api")

	// 6. Handle status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Parse the API response wrapper
		var apiResp tripsAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			log.Error().
				Err(err).
				Str("trip_id", tripID).
				Str("body", string(body)).
				Msg("Failed to parse trips-api response")
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Check if the API returned success=false
		if !apiResp.Success {
			log.Warn().
				Str("trip_id", tripID).
				Str("error", apiResp.Error).
				Msg("trips-api returned success=false")
			return nil, fmt.Errorf("trips-api returned success=false: %s", apiResp.Error)
		}

		// Parse the trip data from the nested data field
		var trip domain.Trip
		if err := json.Unmarshal(apiResp.Data, &trip); err != nil {
			log.Error().
				Err(err).
				Str("trip_id", tripID).
				Str("data", string(apiResp.Data)).
				Msg("Failed to parse trip data")
			return nil, fmt.Errorf("failed to parse trip data: %w", err)
		}

		log.Info().
			Str("trip_id", tripID).
			Int64("driver_id", trip.DriverID).
			Str("status", trip.Status).
			Int("available_seats", trip.AvailableSeats).
			Msg("Successfully retrieved trip from trips-api")

		return &trip, nil

	case http.StatusNotFound:
		log.Warn().Str("trip_id", tripID).Msg("Trip not found in trips-api")
		return nil, domain.ErrTripNotFound.WithDetails(map[string]interface{}{
			"trip_id": tripID,
		})

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("trip_id", tripID).
			Str("body", string(body)).
			Msg("trips-api returned server error")
		return nil, domain.ErrTripsAPIUnavailable.WithDetails(map[string]interface{}{
			"status_code": resp.StatusCode,
			"trip_id":     tripID,
		})

	default:
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("trip_id", tripID).
			Str("body", string(body)).
			Msg("trips-api returned unexpected status code")
		return nil, fmt.Errorf("trips-api returned unexpected status %d: %s", resp.StatusCode, string(body))
	}
}
