package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"payments-api/internal/domain"
)

// User represents basic user information fetched from users-api
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UsersClient defines operations for interacting with users-api
type UsersClient interface {
	// GetUser fetches a user by their numeric ID.
	// Returns domain.ErrUserNotFound if the user does not exist.
	GetUser(ctx context.Context, userID int64) (*User, error)
}

type usersHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUsersClient creates a new HTTP client for users-api
func NewUsersClient(baseURL string) UsersClient {
	return &usersHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type usersAPIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// GetUser fetches a user from users-api by ID.
// Uses X-Internal-Token for service-to-service authentication.
func (c *usersHTTPClient) GetUser(ctx context.Context, userID int64) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Internal-Token", os.Getenv("INTERNAL_API_TOKEN"))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call users-api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var apiResp usersAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		if !apiResp.Success {
			return nil, fmt.Errorf("users-api returned success=false: %s", apiResp.Error)
		}
		var user User
		if err := json.Unmarshal(apiResp.Data, &user); err != nil {
			return nil, fmt.Errorf("failed to parse user data: %w", err)
		}
		return &user, nil

	case http.StatusNotFound:
		return nil, domain.ErrUserNotFound

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return nil, fmt.Errorf("users-api returned %d: service unavailable", resp.StatusCode)

	default:
		return nil, fmt.Errorf("users-api returned unexpected status %d: %s", resp.StatusCode, string(body))
	}
}
