package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserDTO struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type UsersClient interface {
	GetUser(userID int64) (*UserDTO, error)
	ValidateUser(userID int64) error
}

type usersClientImpl struct {
	baseURL string
	client  *http.Client
}

func NewUsersClient(baseURL string) UsersClient {
	return &usersClientImpl{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *usersClientImpl) GetUser(userID int64) (*UserDTO, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Success bool     `json:"success"`
		Data    *UserDTO `json:"data"`
		Error   string   `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("api error: %s", response.Error)
	}

	return response.Data, nil
}

func (c *usersClientImpl) ValidateUser(userID int64) error {
	_, err := c.GetUser(userID)
	return err
}
