package mocks

import (
	"context"
	"time"
)

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	GetFunc         func(ctx context.Context, key string) (string, error)
	SetFunc         func(ctx context.Context, key string, value string, ttl time.Duration) error
	DeleteFunc      func(ctx context.Context, key string) error
	ExistsFunc      func(ctx context.Context, key string) (bool, error)
	IncrementFunc   func(ctx context.Context, key string) (int64, error)
	SetNXFunc       func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	GetWithTTLFunc  func(ctx context.Context, key string) (string, time.Duration, error)
	DeletePatternFunc func(ctx context.Context, pattern string) error
}

// Get calls the mocked GetFunc
func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key)
	}
	return "", nil
}

// Set calls the mocked SetFunc
func (m *MockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, ttl)
	}
	return nil
}

// Delete calls the mocked DeleteFunc
func (m *MockCache) Delete(ctx context.Context, key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key)
	}
	return nil
}

// Exists calls the mocked ExistsFunc
func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, key)
	}
	return false, nil
}

// Increment calls the mocked IncrementFunc
func (m *MockCache) Increment(ctx context.Context, key string) (int64, error) {
	if m.IncrementFunc != nil {
		return m.IncrementFunc(ctx, key)
	}
	return 0, nil
}

// SetNX calls the mocked SetNXFunc
func (m *MockCache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	if m.SetNXFunc != nil {
		return m.SetNXFunc(ctx, key, value, ttl)
	}
	return true, nil
}

// GetWithTTL calls the mocked GetWithTTLFunc
func (m *MockCache) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	if m.GetWithTTLFunc != nil {
		return m.GetWithTTLFunc(ctx, key)
	}
	return "", 0, nil
}

// DeletePattern calls the mocked DeletePatternFunc
func (m *MockCache) DeletePattern(ctx context.Context, pattern string) error {
	if m.DeletePatternFunc != nil {
		return m.DeletePatternFunc(ctx, pattern)
	}
	return nil
}
