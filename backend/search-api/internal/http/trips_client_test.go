package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"search-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTripsClient_GetTrip_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/trips/507f1f77bcf86cd799439011", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		trip := domain.Trip{
			ID:                       primitive.NewObjectID(),
			DriverID:                 1001,
			Status:                   "published",
			AvailableSeats:           3,
			TotalSeats:               4,
			PricePerSeat:             500.0,
			DepartureDatetime:        time.Now().Add(24 * time.Hour),
			EstimatedArrivalDatetime: time.Now().Add(26 * time.Hour),
			Origin: domain.Location{
				City:     "Buenos Aires",
				Province: "Buenos Aires",
				Address:  "Av. Corrientes 1000",
				Coordinates: domain.NewGeoJSONPoint(-34.6037, -58.3816),
			},
			Destination: domain.Location{
				City:     "La Plata",
				Province: "Buenos Aires",
				Address:  "Calle 7 y 50",
				Coordinates: domain.NewGeoJSONPoint(-34.9214, -57.9544),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		resp := StandardResponse{
			Success: true,
			Data:    mustMarshal(trip),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 3,
	})

	// Execute
	trip, err := client.GetTrip(context.Background(), "507f1f77bcf86cd799439011")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, trip)
	assert.Equal(t, int64(1001), trip.DriverID)
	assert.Equal(t, "published", trip.Status)
	assert.Equal(t, 3, trip.AvailableSeats)
	assert.Equal(t, "Buenos Aires", trip.Origin.City)
}

func TestTripsClient_GetTrip_NotFound(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		resp := StandardResponse{
			Success: false,
			Error:   "Trip not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 0, // No retries for 404
	})

	// Execute
	trip, err := client.GetTrip(context.Background(), "nonexistent")

	// Assert
	require.Error(t, err)
	assert.Nil(t, trip)
	assert.Equal(t, domain.ErrTripNotFound, err)
}

func TestTripsClient_GetTrip_ServerError_WithRetry(t *testing.T) {
	attemptCount := 0

	// Create mock server that fails twice then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++

		if attemptCount < 3 {
			// Return 500 for first 2 attempts
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Succeed on 3rd attempt
		trip := domain.Trip{
			ID:             primitive.NewObjectID(),
			DriverID:       1001,
			Status:         "published",
			AvailableSeats: 3,
		}

		resp := StandardResponse{
			Success: true,
			Data:    mustMarshal(trip),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with retries
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:       server.URL,
		Timeout:       5 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 10 * time.Millisecond, // Fast retry for tests
	})

	// Execute
	trip, err := client.GetTrip(context.Background(), "test-trip")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, trip)
	assert.Equal(t, 3, attemptCount, "Should have made 3 attempts (2 retries)")
	assert.Equal(t, "published", trip.Status)
}

func TestTripsClient_GetTrip_Timeout(t *testing.T) {
	// Create mock server with slow response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		json.NewEncoder(w).Encode(StandardResponse{Success: true})
	}))
	defer server.Close()

	// Create client with short timeout
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    100 * time.Millisecond,
		MaxRetries: 0,
	})

	// Execute
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	trip, err := client.GetTrip(ctx, "test-trip")

	// Assert
	require.Error(t, err)
	assert.Nil(t, trip)
}

func TestTripsClient_GetTrip_InvalidJSON(t *testing.T) {
	// Create mock server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	// Create client
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 0,
	})

	// Execute
	trip, err := client.GetTrip(context.Background(), "test-trip")

	// Assert
	require.Error(t, err)
	assert.Nil(t, trip)
	assert.True(t, domain.IsNotFoundError(err) || err.Error() != "")
}

func TestTripsClient_CircuitBreaker(t *testing.T) {
	attemptCount := 0

	// Create mock server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client with circuit breaker (3 failures) and reduced retries
	client := NewTripsClient(HTTPClientConfig{
		BaseURL:        server.URL,
		Timeout:        5 * time.Second,
		MaxRetries:     1, // 1 retry to reduce test time
		RetryWaitTime:  10 * time.Millisecond,
		CircuitBreaker: NewCircuitBreaker(3, 10*time.Second),
	})

	// Make 5 requests - circuit should open after 3 failures
	var lastErr error
	for i := 0; i < 5; i++ {
		_, err := client.GetTrip(context.Background(), "test-trip")
		require.Error(t, err)
		lastErr = err

		if i >= 3 {
			// After 3 failures, circuit should be open and requests fail immediately
			assert.Equal(t, domain.ErrServiceUnavailable, err, "Circuit should be open")
		}
	}

	// Last error should be circuit breaker error
	assert.Equal(t, domain.ErrServiceUnavailable, lastErr)
}

// Helper function to marshal data
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
