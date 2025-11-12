// internal/cache/redis.go
package cache

import (
	"context"
	"fmt"
	"time"
)

// RedisCache implementa la interface Cache usando Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache crea una nueva instancia de RedisCache
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Verificar conexión con ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get obtiene un valor del cache por su key
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key no existe, retornar string vacío (no es error)
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error getting key %s: %w", key, err)
	}
	return val, nil
}

// Set guarda un valor en el cache con un TTL
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}
	return nil
}

// Delete elimina una key del cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting key %s: %w", key, err)
	}
	return nil
}

// Exists verifica si una key existe en el cache
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking existence of key %s: %w", key, err)
	}
	return result > 0, nil
}

// Close cierra la conexión con Redis
func (r *RedisCache) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// GetWithTTL obtiene un valor y su TTL restante
func (r *RedisCache) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	pipe := r.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return "", 0, fmt.Errorf("error in pipeline: %w", err)
	}

	val, err := getCmd.Result()
	if err == redis.Nil {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, fmt.Errorf("error getting key: %w", err)
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return val, 0, fmt.Errorf("error getting TTL: %w", err)
	}

	return val, ttl, nil
}

// SetNX (Set if Not eXists) - útil para locks o cache único
func (r *RedisCache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	result, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("error in SetNX for key %s: %w", key, err)
	}
	return result, nil
}

// Increment incrementa un valor numérico (útil para contadores)
func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	result, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing key %s: %w", key, err)
	}
	return result, nil
}

// DeletePattern elimina todas las keys que coincidan con un patrón
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	var cursor uint64
	var deleted int64

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return deleted, fmt.Errorf("error scanning keys: %w", err)
		}

		if len(keys) > 0 {
			n, err := r.client.Del(ctx, keys...).Result()
			if err != nil {
				return deleted, fmt.Errorf("error deleting keys: %w", err)
			}
			deleted += n
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return deleted, nil
}

// Ping verifica la conexión con Redis
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Stats retorna estadísticas del cliente Redis
func (r *RedisCache) Stats() *redis.PoolStats {
	return r.client.PoolStats()
}
