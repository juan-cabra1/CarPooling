package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthController handles health check endpoints
type HealthController struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthController creates a new HealthController
func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{
		db:        db,
		startTime: time.Now(),
	}
}

// Health returns the service health status
// @Summary Health check
// @Description Returns the health status of the payments-api service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthController) Health(c *gin.Context) {
	// Check database connection
	dbStatus := "healthy"
	sqlDB, err := h.db.DB()
	if err != nil {
		dbStatus = "unhealthy: " + err.Error()
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy: " + err.Error()
	}

	uptime := time.Since(h.startTime)

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "payments-api",
		"version": "1.0.0",
		"uptime":  uptime.String(),
		"checks": gin.H{
			"database": dbStatus,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready returns whether the service is ready to accept traffic
// @Summary Readiness check
// @Description Returns whether the service is ready to accept traffic
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthController) Ready(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":  false,
			"reason": "database connection error: " + err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":  false,
			"reason": "database ping failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready": true,
	})
}
