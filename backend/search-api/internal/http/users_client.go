package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"search-api/internal/domain"

	"github.com/rs/zerolog/log"
)

// UsersClient defines the interface for communicating with users-api
type UsersClient interface {
	GetUser(ctx context.Context, userID int64) (*domain.User, error)
}

// usersHTTPClient implements UsersClient using HTTP
type usersHTTPClient struct {
	baseURL        string
	client         *http.Client
	maxRetries     int
	retryWaitTime  time.Duration
	circuitBreaker *CircuitBreaker
}

// NewUsersClient creates a new UsersClient with the given configuration
func NewUsersClient(config HTTPClientConfig) UsersClient {
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

	return &usersHTTPClient{
		baseURL:        config.BaseURL,
		client:         CreateHTTPClient(config.Timeout),
		maxRetries:     config.MaxRetries,
		retryWaitTime:  config.RetryWaitTime,
		circuitBreaker: config.CircuitBreaker,
	}
}

// GetUser fetches user details from users-api
// Endpoint: GET /users/:id
// Returns: User DTO with profile and rating information
func (c *usersHTTPClient) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	log.Debug().
		Int64("user_id", userID).
		Str("url", url).
		Msg("Fetching user from users-api")

	// Execute with circuit breaker
	var user *domain.User
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

		// Handle 404 specifically for users
		if resp.StatusCode == http.StatusNotFound {
			return domain.ErrUserNotFound
		}

		// Parse standard response
		var userData domain.User
		if err := ParseStandardResponse(resp, &userData); err != nil {
			return err
		}

		user = &userData
		return nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("Failed to fetch user from users-api")
		return nil, err
	}

	log.Info().
		Int64("user_id", userID).
		Str("name", user.Name).
		Float64("rating", user.AverageRatingAsDriver).
		Int("total_trips", user.TotalTripsAsDriver).
		Msg("Successfully fetched user from users-api")

	return user, nil
}
