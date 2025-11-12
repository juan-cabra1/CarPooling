// internal/cache/cache.go
package cache

import (
	"context"
	"time"
)

// Cache define las operaciones básicas de cache
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}
