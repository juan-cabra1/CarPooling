package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type windowEntry struct {
	mu        sync.Mutex
	count     int
	windowEnd time.Time
}

var ipStore sync.Map

func init() {
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
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, _ := ipStore.LoadOrStore(ip, &windowEntry{})
		entry := val.(*windowEntry)

		entry.mu.Lock()
		if now.After(entry.windowEnd) {
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
