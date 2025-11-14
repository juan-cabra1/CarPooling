package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Success(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogger_WithQueryParams(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/search", func(c *gin.Context) {
		c.JSON(200, gin.H{"results": []string{"trip1", "trip2"}})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?origin=Buenos+Aires&destination=Cordoba", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogger_WithResultsCount(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/search", func(c *gin.Context) {
		// Simulate setting results_count in context
		c.Set("results_count", 15)
		c.JSON(200, gin.H{"results": []string{"trip1", "trip2"}})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogger_WithSource(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/search", func(c *gin.Context) {
		// Simulate setting source in context
		c.Set("source", "solr")
		c.Set("results_count", 10)
		c.JSON(200, gin.H{"results": []string{"trip1"}})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogger_Error4xx(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "bad request"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestLogger_Error5xx(t *testing.T) {
	router := gin.New()
	router.Use(Logger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "internal error"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}

func TestLogger_DifferentMethods(t *testing.T) {
	tests := []struct {
		method     string
		statusCode int
	}{
		{"GET", 200},
		{"POST", 201},
		{"PUT", 200},
		{"DELETE", 204},
		{"PATCH", 200},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			router := gin.New()
			router.Use(Logger())
			router.Handle(tt.method, "/test", func(c *gin.Context) {
				c.Status(tt.statusCode)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
