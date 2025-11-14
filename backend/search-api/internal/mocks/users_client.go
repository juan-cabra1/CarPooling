package mocks

import (
	"context"

	"github.com/juan-cabra1/CarPooling/backend/search-api/internal/domain"
)

// MockUsersClient is a mock implementation of the UsersClient interface
type MockUsersClient struct {
	GetUserFunc func(ctx context.Context, userID string) (*domain.User, error)
}

// GetUser calls the mocked GetUserFunc
func (m *MockUsersClient) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, userID)
	}
	return nil, nil
}
