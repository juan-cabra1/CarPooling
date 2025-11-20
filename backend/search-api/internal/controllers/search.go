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
		SearchText: c.Query("q"),
		SortBy:     c.DefaultQuery("sort_by", "earliest"),
		SortOrder:  c.DefaultQuery("sort_order", "asc"),
	}

	// Parse Origin Location
	query.Origin = parseLocation(c, "origin")

	// Parse Destination Location
	query.Destination = parseLocation(c, "destination")

	// Parse Origin Radius (optional, for geospatial search)
	if radiusStr := c.Query("origin_radius"); radiusStr != "" {
		if val, err := strconv.Atoi(radiusStr); err == nil {
			query.OriginRadius = val
		}
	}

	// Parse Destination Radius (optional, for geospatial search)
	if radiusStr := c.Query("destination_radius"); radiusStr != "" {
		if val, err := strconv.Atoi(radiusStr); err == nil {
			query.DestinationRadius = val
		}
	}

	// Parse Departure Date (exact date)
	if dateStr := c.Query("departure_date"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			query.DepartureDate = &t
		} else if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			query.DepartureDate = &t
		}
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

	// Parse pagination
	query.Page = parseInt(c.DefaultQuery("page", "1"))
	query.Limit = parseInt(c.DefaultQuery("limit", "20"))

	// Set defaults and validate
	query.SetDefaults()
	if err := query.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_QUERY",
				"message": err.Error(),
			},
		})
		return
	}

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
// DEPRECATED: Use SearchTrips with origin coordinates instead
func (sc *SearchController) SearchByLocation(c *gin.Context) {
	// Parse geospatial parameters
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius_km")

	if latStr == "" || lngStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_GEO_PARAMS",
				"message": "lat, lng, and radius_km are all required",
			},
		})
		return
	}

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

	// Build query using new structure
	query := &domain.SearchQuery{
		Origin: &domain.Location{
			Coordinates: domain.NewGeoJSONPoint(lat, lng),
		},
		OriginRadius: radiusKm,
	}

	// Parse additional filters
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

	query.Page = parseInt(c.DefaultQuery("page", "1"))
	query.Limit = parseInt(c.DefaultQuery("limit", "20"))

	query.SetDefaults()
	if err := query.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_QUERY",
				"message": err.Error(),
			},
		})
		return
	}

	// Call unified search service
	results, err := sc.searchService.SearchTrips(c.Request.Context(), query)
	if err != nil {
		c.Error(err)
		return
	}

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

	suggestions, err := sc.searchService.GetAutocomplete(c.Request.Context(), query, limit)
	if err != nil {
		c.Error(err)
		return
	}

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

	routes, err := sc.searchService.GetPopularRoutes(c.Request.Context(), limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"routes": routes,
		},
	})
}

// Helper functions

// parseLocation parses a Location from query parameters
// Supports formats (with backward compatibility):
// - NEW: ?origin_city=Córdoba&origin_province=Córdoba
// - NEW: ?origin_city=Córdoba&origin_province=Córdoba&origin_lat=-31.4&origin_lng=-64.2
// - OLD (deprecated): ?originCity=Córdoba&originProvince=Córdoba&originLat=-31.4&originLng=-64.2
func parseLocation(c *gin.Context, prefix string) *domain.Location {
	// Parse city and province - try new format first (snake_case)
	city := c.Query(prefix + "_city")
	province := c.Query(prefix + "_province")

	// Fallback to old format (camelCase) for backward compatibility
	if city == "" {
		city = c.Query(prefix + "City")
	}
	if province == "" {
		province = c.Query(prefix + "Province")
	}

	// Parse coordinates - try new format first
	latStr := c.Query(prefix + "_lat")
	lngStr := c.Query(prefix + "_lng")

	// Fallback to old format for coordinates
	if latStr == "" {
		latStr = c.Query(prefix + "Lat")
	}
	if lngStr == "" {
		lngStr = c.Query(prefix + "Lng")
	}

	// Parse coordinate values if provided
	var hasCoordinates bool
	var lat, lng float64
	if latStr != "" && lngStr != "" {
		parsedLat, err1 := strconv.ParseFloat(latStr, 64)
		parsedLng, err2 := strconv.ParseFloat(lngStr, 64)

		if err1 == nil && err2 == nil {
			lat = parsedLat
			lng = parsedLng
			hasCoordinates = true
		}
	}

	// Return nil only if BOTH city and coordinates are missing
	if city == "" && !hasCoordinates {
		return nil
	}

	// Create location with available data
	location := &domain.Location{
		City:        city,
		Province:    province,
		Address:     "", // Not used in search
		Coordinates: domain.GeoJSONPoint{Type: "Point", Coordinates: []float64{}}, // Empty by default
	}

	// Set coordinates if they were successfully parsed
	if hasCoordinates {
		location.Coordinates = domain.NewGeoJSONPoint(lat, lng)
	}

	return location
}

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
