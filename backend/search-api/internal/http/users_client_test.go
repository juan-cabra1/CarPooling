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
)

func TestUsersClient_GetUser_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/1001", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		user := domain.User{
			ID:                       1001,
			Name:                     "Juan Pérez",
			Email:                    "juan@example.com",
			PhoneNumber:              "+54911234567",
			PhotoURL:                 "https://example.com/photo.jpg",
			IsVerified:               true,
			TotalTripsAsDriver:       25,
			TotalTripsAsPassenger:    10,
			AverageRatingAsDriver:    4.8,
			TotalRatingsAsDriver:     20,
			AverageRatingAsPassenger: 4.9,
			TotalRatingsAsPassenger:  8,
			CreatedAt:                time.Now().Add(-365 * 24 * time.Hour),
			UpdatedAt:                time.Now(),
		}

		resp := StandardResponse{
			Success: true,
			Data:    mustMarshal(user),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 3,
	})

	// Execute
	user, err := client.GetUser(context.Background(), 1001)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(1001), user.ID)
	assert.Equal(t, "Juan Pérez", user.Name)
	assert.Equal(t, "juan@example.com", user.Email)
	assert.Equal(t, 25, user.TotalTripsAsDriver)
	assert.Equal(t, 4.8, user.AverageRatingAsDriver)
	assert.True(t, user.IsVerified)
}

func TestUsersClient_GetUser_NotFound(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		resp := StandardResponse{
			Success: false,
			Error:   "User not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 0, // No retries for 404
	})

	// Execute
	user, err := client.GetUser(context.Background(), 9999)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrUserNotFound, err)
}

func TestUsersClient_GetUser_ServerError_WithRetry(t *testing.T) {
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
		user := domain.User{
			ID:                    1001,
			Name:                  "Juan Pérez",
			Email:                 "juan@example.com",
			TotalTripsAsDriver:    25,
			AverageRatingAsDriver: 4.8,
		}

		resp := StandardResponse{
			Success: true,
			Data:    mustMarshal(user),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with retries
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:       server.URL,
		Timeout:       5 * time.Second,
		MaxRetries:    3,
		RetryWaitTime: 10 * time.Millisecond, // Fast retry for tests
	})

	// Execute
	user, err := client.GetUser(context.Background(), 1001)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 3, attemptCount, "Should have made 3 attempts (2 retries)")
	assert.Equal(t, "Juan Pérez", user.Name)
}

func TestUsersClient_GetUser_Timeout(t *testing.T) {
	// Create mock server with slow response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		json.NewEncoder(w).Encode(StandardResponse{Success: true})
	}))
	defer server.Close()

	// Create client with short timeout
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    100 * time.Millisecond,
		MaxRetries: 0,
	})

	// Execute
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	user, err := client.GetUser(ctx, 1001)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
}

func TestUsersClient_GetUser_ToDriver(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := domain.User{
			ID:                    1001,
			Name:                  "Juan Pérez",
			Email:                 "juan@example.com",
			PhotoURL:              "https://example.com/photo.jpg",
			TotalTripsAsDriver:    25,
			AverageRatingAsDriver: 4.8,
		}

		resp := StandardResponse{
			Success: true,
			Data:    mustMarshal(user),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 3,
	})

	// Execute
	user, err := client.GetUser(context.Background(), 1001)
	require.NoError(t, err)

	// Convert to Driver
	driver := user.ToDriver()

	// Assert Driver fields
	assert.Equal(t, int64(1001), driver.ID)
	assert.Equal(t, "Juan Pérez", driver.Name)
	assert.Equal(t, "juan@example.com", driver.Email)
	assert.Equal(t, "https://example.com/photo.jpg", driver.PhotoURL)
	assert.Equal(t, 4.8, driver.Rating)
	assert.Equal(t, 25, driver.TotalTrips)
}

func TestUsersClient_CircuitBreaker(t *testing.T) {
	attemptCount := 0

	// Create mock server that always fails
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create client with circuit breaker (3 failures) and reduced retries
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:        server.URL,
		Timeout:        5 * time.Second,
		MaxRetries:     1, // 1 retry to reduce test time
		RetryWaitTime:  10 * time.Millisecond,
		CircuitBreaker: NewCircuitBreaker(3, 10*time.Second),
	})

	// Make 5 requests - circuit should open after 3 failures
	var lastErr error
	for i := 0; i < 5; i++ {
		_, err := client.GetUser(context.Background(), 1001)
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

func TestUsersClient_GetUser_UnauthorizedError(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		resp := StandardResponse{
			Success: false,
			Error:   "Unauthorized",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 0,
	})

	// Execute
	user, err := client.GetUser(context.Background(), 1001)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestUsersClient_GetUser_InvalidJSON(t *testing.T) {
	// Create mock server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	// Create client
	client := NewUsersClient(HTTPClientConfig{
		BaseURL:    server.URL,
		Timeout:    5 * time.Second,
		MaxRetries: 0,
	})

	// Execute
	user, err := client.GetUser(context.Background(), 1001)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
}
