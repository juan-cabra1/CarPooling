package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check requests
// This controller provides a simple endpoint to verify the service is running
type HealthController struct {
	serviceName string
	servicePort string
}

// NewHealthController creates a new HealthController instance
//
// Parameters:
//   - serviceName: Name of the service (e.g., "bookings-api")
//   - servicePort: Port the service is running on (e.g., "8003")
//
// Returns:
//   - *HealthController: Initialized health controller
func NewHealthController(serviceName, servicePort string) *HealthController {
	return &HealthController{
		serviceName: serviceName,
		servicePort: servicePort,
	}
}

// HealthCheck handles GET /health requests
// Returns a simple JSON response indicating the service is operational
//
// This endpoint is used by:
//   - Load balancers for health checks
//   - Monitoring systems (Prometheus, Datadog, etc.)
//   - Docker/Kubernetes health probes
//   - Manual service verification
//
// Response format:
//
//	{
//	  "status": "ok",
//	  "service": "bookings-api",
//	  "port": "8003"
//	}
//
// HTTP Status: 200 OK
func (h *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": h.serviceName,
		"port":    h.servicePort,
	})
}
