package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"search-api/internal/domain"

	"github.com/rs/zerolog/log"
)

// StandardResponse represents the standard API response wrapper used by trips-api and users-api
type StandardResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// HTTPClientConfig holds configuration for HTTP clients
type HTTPClientConfig struct {
	BaseURL        string
	Timeout        time.Duration
	MaxRetries     int
	RetryWaitTime  time.Duration
	CircuitBreaker *CircuitBreaker
}

// CircuitBreaker implements a simple circuit breaker pattern
type CircuitBreaker struct {
	MaxFailures       int
	ResetTimeout      time.Duration
	consecutiveFails  int
	lastFailureTime   time.Time
	isOpen            bool
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		MaxFailures:  maxFailures,
		ResetTimeout: resetTimeout,
	}
}

// Call executes the function if circuit is closed, returns error if open
func (cb *CircuitBreaker) Call(fn func() error) error {
	// Check if circuit should be reset
	if cb.isOpen && time.Since(cb.lastFailureTime) > cb.ResetTimeout {
		log.Info().Msg("Circuit breaker reset - trying half-open state")
		cb.isOpen = false
		cb.consecutiveFails = 0
	}

	// Circuit is open - fail fast
	if cb.isOpen {
		return domain.ErrServiceUnavailable
	}

	// Execute function
	err := fn()
	if err != nil {
		cb.consecutiveFails++
		cb.lastFailureTime = time.Now()

		// Open circuit if max failures reached
		if cb.consecutiveFails >= cb.MaxFailures {
			cb.isOpen = true
			log.Warn().
				Int("consecutive_failures", cb.consecutiveFails).
				Msg("Circuit breaker opened - service marked as unavailable")
		}
		return err
	}

	// Success - reset counter
	if cb.consecutiveFails > 0 {
		log.Info().Msg("Circuit breaker reset - service recovered")
	}
	cb.consecutiveFails = 0
	return nil
}

// DoRequestWithRetry executes an HTTP request with retry logic and logging
func DoRequestWithRetry(
	ctx context.Context,
	client *http.Client,
	req *http.Request,
	maxRetries int,
	retryWaitTime time.Duration,
) (*http.Response, error) {
	var resp *http.Response
	var err error

	startTime := time.Now()
	serviceName := req.URL.Host

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			waitDuration := retryWaitTime * time.Duration(1<<uint(attempt-1))
			log.Warn().
				Str("service", serviceName).
				Str("url", req.URL.String()).
				Int("attempt", attempt).
				Dur("wait_duration", waitDuration).
				Msg("Retrying HTTP request")

			select {
			case <-time.After(waitDuration):
			case <-ctx.Done():
				return nil, domain.WrapError(domain.ErrTimeout, "context canceled during retry")
			}
		}

		// Execute request
		resp, err = client.Do(req.WithContext(ctx))

		// Success - return immediately
		if err == nil && resp.StatusCode < 500 {
			duration := time.Since(startTime)
			log.Info().
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Int("status_code", resp.StatusCode).
				Int64("duration_ms", duration.Milliseconds()).
				Int("attempt", attempt+1).
				Msg("HTTP request completed")
			return resp, nil
		}

		// Log error for retry
		if err != nil {
			log.Error().
				Err(err).
				Str("service", serviceName).
				Str("url", req.URL.String()).
				Int("attempt", attempt+1).
				Msg("HTTP request failed - network error")
		} else if resp.StatusCode >= 500 {
			log.Error().
				Str("service", serviceName).
				Str("url", req.URL.String()).
				Int("status_code", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("HTTP request failed - server error")
			resp.Body.Close() // Close body before retry
		}

		// Don't retry if we've reached max attempts
		if attempt == maxRetries {
			break
		}
	}

	// All retries exhausted
	duration := time.Since(startTime)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", serviceName).
			Str("url", req.URL.String()).
			Int64("duration_ms", duration.Milliseconds()).
			Int("max_retries", maxRetries).
			Msg("HTTP request failed after all retries")
		return nil, domain.WrapError(domain.ErrServiceUnavailable, "all retries exhausted")
	}

	return resp, nil
}

// ParseStandardResponse parses the standard API response wrapper and extracts data
func ParseStandardResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.WrapError(domain.ErrInvalidResponse, "failed to read response body")
	}

	// Handle error status codes
	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrTripNotFound // Will be overridden by specific clients
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.ErrUnauthorized
	}

	if resp.StatusCode >= 500 {
		return domain.ErrServiceUnavailable
	}

	if resp.StatusCode != http.StatusOK {
		return domain.NewAppError(
			"UNEXPECTED_STATUS",
			fmt.Sprintf("unexpected status code: %d", resp.StatusCode),
			string(body),
		)
	}

	// Parse standard response wrapper
	var stdResp StandardResponse
	if err := json.Unmarshal(body, &stdResp); err != nil {
		log.Error().
			Err(err).
			Str("body", string(body)).
			Msg("Failed to parse standard response")
		return domain.WrapError(domain.ErrInvalidResponse, "failed to parse JSON response")
	}

	// Check if response indicates failure
	if !stdResp.Success {
		return domain.NewAppError(
			"API_ERROR",
			fmt.Sprintf("API returned error: %s", stdResp.Error),
			stdResp.Error,
		)
	}

	// Unmarshal data into target
	if err := json.Unmarshal(stdResp.Data, target); err != nil {
		log.Error().
			Err(err).
			Str("data", string(stdResp.Data)).
			Msg("Failed to unmarshal response data")
		return domain.WrapError(domain.ErrInvalidResponse, "failed to unmarshal data field")
	}

	return nil
}

// CreateHTTPClient creates a configured HTTP client with timeout
func CreateHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}
