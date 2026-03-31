package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// windowEntry representa el estado de un IP en la ventana de tiempo actual
type windowEntry struct {
	mu        sync.Mutex
	count     int
	windowEnd time.Time
}

// ipStore mantiene los contadores por IP usando ventana fija de 1 minuto
var ipStore sync.Map

func init() {
	// Limpieza periódica de entradas expiradas para evitar memory leak
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			now := time.Now()
			ipStore.Range(func(key, value any) bool {
				entry := value.(*windowEntry)
				entry.mu.Lock()
				expired := now.After(entry.windowEnd)
				entry.mu.Unlock()
				if expired {
					ipStore.Delete(key)
				}
				return true
			})
		}
	}()
}

// RateLimitMiddleware limita requests por IP usando ventana fija de 1 minuto.
// requestsPerMinute: máximo de requests permitidos por IP por minuto.
//
// Ejemplo de uso:
//
//	router.POST("/login", middleware.RateLimitMiddleware(10), handler) // 10 req/min
//	protected.Use(middleware.RateLimitMiddleware(120))                 // 120 req/min para rutas autenticadas
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, _ := ipStore.LoadOrStore(ip, &windowEntry{})
		entry := val.(*windowEntry)

		entry.mu.Lock()
		if now.After(entry.windowEnd) {
			// Nueva ventana
			entry.count = 1
			entry.windowEnd = now.Add(time.Minute)
		} else {
			entry.count++
		}
		count := entry.count
		entry.mu.Unlock()

		if count > requestsPerMinute {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "demasiadas solicitudes, intentá de nuevo en un minuto",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
