package controllers

import (
	"net/http"
	"search-api/internal/domain"
	"search-api/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SearchController handles all search-related endpoints
type SearchController struct {
	searchService service.SearchService
}

// NewSearchController creates a new SearchController instance
func NewSearchController(searchService service.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

// SearchTrips handles GET /api/v1/search/trips
func (sc *SearchController) SearchTrips(c *gin.Context) {
	// Build query from query parameters
	query := &domain.SearchQuery{
		OriginCity:      c.Query("origin_city"),
		DestinationCity: c.Query("destination_city"),
		SearchText:      c.Query("q"),
		SortBy:          c.DefaultQuery("sort_by", "popularity"),
	}

	// Parse numeric filters
	if minSeats := c.Query("min_seats"); minSeats != "" {
		if val, err := strconv.Atoi(minSeats); err == nil {
			query.MinSeats = val
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query.MaxPrice = val
		}
	}
	if minRating := c.Query("min_driver_rating"); minRating != "" {
		if val, err := strconv.ParseFloat(minRating, 64); err == nil {
			query.MinDriverRating = val
		}
	}

	// Parse boolean filters (use pointers for true/false/unset)
	query.PetsAllowed = parseBoolPtr(c, "pets_allowed")
	query.SmokingAllowed = parseBoolPtr(c, "smoking_allowed")
	query.MusicAllowed = parseBoolPtr(c, "music_allowed")

	// Parse date filters (ISO8601)
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			query.DateFrom = t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			query.DateTo = t
		}
	}

	// Parse pagination
	query.Page = parseInt(c.DefaultQuery("page", "1"))
	query.Limit = parseInt(c.DefaultQuery("limit", "20"))

	// Call service
	results, err := sc.searchService.SearchTrips(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"trips":       results.Trips,
			"total":       results.Total,
			"page":        results.Page,
			"limit":       results.Limit,
			"total_pages": results.TotalPages,
		},
	})
}

// SearchByLocation handles GET /api/v1/search/location
func (sc *SearchController) SearchByLocation(c *gin.Context) {
	// Parse required geospatial parameters
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius_km")

	// Validate all 3 are present
	if latStr == "" || lngStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_GEO_PARAMS",
				"message": "lat, lng, and radius_km are all required for location search",
			},
		})
		return
	}

	// Parse coordinates
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_LATITUDE",
				"message": "Invalid latitude value",
			},
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_LONGITUDE",
				"message": "Invalid longitude value",
			},
		})
		return
	}

	radiusKm, err := strconv.Atoi(radiusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_RADIUS",
				"message": "Invalid radius_km value",
			},
		})
		return
	}

	// Build additional filters
	filters := make(map[string]interface{})
	if minSeats := c.Query("min_seats"); minSeats != "" {
		if val, err := strconv.Atoi(minSeats); err == nil {
			filters["min_seats"] = val
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = val
		}
	}

	// Call service
	results, err := sc.searchService.SearchByLocation(c.Request.Context(), lat, lng, radiusKm, filters)
	if err != nil {
		c.Error(err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"trips":       results.Trips,
			"total":       results.Total,
			"page":        results.Page,
			"limit":       results.Limit,
			"total_pages": results.TotalPages,
		},
	})
}

// GetTrip handles GET /api/v1/trips/:id
func (sc *SearchController) GetTrip(c *gin.Context) {
	tripID := c.Param("id")
	if tripID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MISSING_TRIP_ID",
				"message": "Trip ID is required",
			},
		})
		return
	}

	// Call service
	trip, err := sc.searchService.GetTrip(c.Request.Context(), tripID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TRIP_NOT_FOUND",
					"message": "Trip not found",
				},
			})
			return
		}
		c.Error(err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"trip": trip,
		},
	})
}

// GetAutocomplete handles GET /api/v1/search/autocomplete
func (sc *SearchController) GetAutocomplete(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "QUERY_TOO_SHORT",
				"message": "Query must be at least 2 characters",
			},
		})
		return
	}

	limit := parseInt(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}

	// Call service
	suggestions, err := sc.searchService.GetAutocomplete(c.Request.Context(), query, limit)
	if err != nil {
		c.Error(err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"suggestions": suggestions,
		},
	})
}

// GetPopularRoutes handles GET /api/v1/search/popular-routes
func (sc *SearchController) GetPopularRoutes(c *gin.Context) {
	limit := parseInt(c.DefaultQuery("limit", "10"))
	if limit > 50 {
		limit = 50
	}

	// Call service
	routes, err := sc.searchService.GetPopularRoutes(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"routes": routes,
		},
	})
}

// Helper functions

// parseBoolPtr parses a boolean query parameter and returns a pointer
// Returns nil if parameter is not present
func parseBoolPtr(c *gin.Context, key string) *bool {
	if val, exists := c.GetQuery(key); exists {
		boolVal := val == "true" || val == "1"
		return &boolVal
	}
	return nil
}

// parseInt safely parses an integer with default fallback
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
