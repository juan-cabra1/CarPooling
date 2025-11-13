package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"search-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestErrorHandler_NoErrors(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestErrorHandler_AppError_NotFound(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(domain.ErrTripNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
	assert.Contains(t, w.Body.String(), `"code":"TRIP_NOT_FOUND"`)
	assert.Contains(t, w.Body.String(), `"message":"Trip not found in trips-api"`)
}

func TestErrorHandler_AppError_InvalidQuery(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(domain.ErrInvalidQuery)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"INVALID_QUERY"`)
}

func TestErrorHandler_AppError_Unauthorized(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(domain.ErrUnauthorized)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"UNAUTHORIZED"`)
}

func TestErrorHandler_AppError_ServiceUnavailable(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(domain.ErrSolrUnavailable)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"SOLR_UNAVAILABLE"`)
}

func TestErrorHandler_GenericError(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("generic error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"INTERNAL_ERROR"`)
	assert.Contains(t, w.Body.String(), `"message":"An internal error occurred"`)
}

func TestErrorHandler_AppErrorWithDetails(t *testing.T) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		customErr := domain.NewAppError(
			"CUSTOM_ERROR",
			"Custom error message",
			gin.H{"field": "origin", "value": "invalid"},
		)
		c.Error(customErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"CUSTOM_ERROR"`)
	assert.Contains(t, w.Body.String(), `"details":{"field":"origin","value":"invalid"}`)
}

func TestMapErrorCodeToHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		errorCode  string
		statusCode int
	}{
		{"TripNotFound", "TRIP_NOT_FOUND", http.StatusNotFound},
		{"UserNotFound", "USER_NOT_FOUND", http.StatusNotFound},
		{"SearchTripNotFound", "SEARCH_TRIP_NOT_FOUND", http.StatusNotFound},
		{"InvalidQuery", "INVALID_QUERY", http.StatusBadRequest},
		{"InvalidGeoCoords", "INVALID_GEO_COORDS", http.StatusBadRequest},
		{"InvalidInput", "INVALID_INPUT", http.StatusBadRequest},
		{"Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized},
		{"SolrUnavailable", "SOLR_UNAVAILABLE", http.StatusServiceUnavailable},
		{"ServiceUnavailable", "SERVICE_UNAVAILABLE", http.StatusServiceUnavailable},
		{"UnknownError", "UNKNOWN_ERROR", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := mapErrorCodeToHTTPStatus(tt.errorCode)
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
