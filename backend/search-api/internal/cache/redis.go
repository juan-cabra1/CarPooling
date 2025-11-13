// internal/cache/redis.go
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// RedisCache implementa la interface Cache usando Memcache
// Mantiene el nombre RedisCache por compatibilidad con código existente
type RedisCache struct {
	client *memcache.Client
}

// NewRedisCache crea una nueva instancia de RedisCache (usando Memcache internamente)
// Los parámetros password y db son ignorados ya que Memcache no los usa
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	// Memcache acepta múltiples servidores, pero pasamos uno solo
	client := memcache.New(addr)

	// Configurar timeouts
	client.Timeout = 3 * time.Second
	client.MaxIdleConns = 10

	// Verificar conexión con ping
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to Memcache: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get obtiene un valor del cache por su key
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	item, err := r.client.Get(key)
	if err == memcache.ErrCacheMiss {
		// Key no existe, retornar string vacío (no es error)
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error getting key %s: %w", key, err)
	}
	return string(item.Value), nil
}

// Set guarda un valor en el cache con un TTL
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(ttl.Seconds()),
	}

	err := r.client.Set(item)
	if err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}
	return nil
}

// Delete elimina una key del cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Delete(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return fmt.Errorf("error deleting key %s: %w", key, err)
	}
	// Si no existe, no es error
	return nil
}

// Exists verifica si una key existe en el cache
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking existence of key %s: %w", key, err)
	}
	return true, nil
}

// Close cierra la conexión con Memcache
// Memcache no requiere close explícito, pero mantenemos el método por compatibilidad
func (r *RedisCache) Close() error {
	// Memcache cierra conexiones automáticamente
	return nil
}

// GetWithTTL obtiene un valor y su TTL restante
// NOTA: Memcache no soporta obtener el TTL restante, por lo que retornamos 0
func (r *RedisCache) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	item, err := r.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return "", 0, nil
	}
	if err != nil {
		return "", 0, fmt.Errorf("error getting key: %w", err)
	}

	// Memcache no expone el TTL restante, retornamos 0
	return string(item.Value), 0, nil
}

// SetNX (Set if Not eXists) - útil para locks o cache único
func (r *RedisCache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(ttl.Seconds()),
	}

	err := r.client.Add(item)
	if err == memcache.ErrNotStored {
		// La key ya existe
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error in SetNX for key %s: %w", key, err)
	}
	return true, nil
}

// Increment incrementa un valor numérico (útil para contadores)
func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	// Memcache Increment requiere que la key ya exista
	// Intentamos incrementar, si no existe la creamos con valor 1
	newValue, err := r.client.Increment(key, 1)
	if err == memcache.ErrCacheMiss {
		// Key no existe, la creamos con valor 1
		item := &memcache.Item{
			Key:        key,
			Value:      []byte("1"),
			Expiration: 0, // Sin expiración
		}
		if err := r.client.Add(item); err != nil {
			return 0, fmt.Errorf("error creating counter key %s: %w", key, err)
		}
		return 1, nil
	}
	if err != nil {
		return 0, fmt.Errorf("error incrementing key %s: %w", key, err)
	}
	return int64(newValue), nil
}

// DeletePattern elimina todas las keys que coincidan con un patrón
// NOTA: Memcache no soporta búsqueda por patrón, esta funcionalidad está limitada
// Para una implementación completa, necesitarías mantener un índice de keys
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	// Memcache no soporta SCAN o búsqueda por patrón
	// Esta es una limitación conocida de Memcache
	// Retornamos error indicando que no está soportado
	return 0, fmt.Errorf("DeletePattern not supported by Memcache (pattern: %s)", pattern)
}

// Ping verifica la conexión con Memcache
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping()
}

// Stats retorna estadísticas del servidor Memcache
// Retorna nil ya que el formato de stats es diferente entre Redis y Memcache
func (r *RedisCache) Stats() interface{} {
	// Memcache tiene su propio formato de stats
	// Para mantener compatibilidad, retornamos nil
	return nil
}
