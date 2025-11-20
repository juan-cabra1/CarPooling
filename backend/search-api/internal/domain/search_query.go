package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// SearchQuery contains all filter parameters for searching trips
type SearchQuery struct {
	//locations
	Origin            *Location `json:"origin,omitempty"`
	Destination       *Location `json:"destination,omitempty"`
	OriginRadius      int       `json:"origin_radius,omitempty"`      // in kilometers
	DestinationRadius int       `json:"destination_radius,omitempty"` // in kilometers

	// Date filter
	DepartureDate *time.Time `json:"departure_date,omitempty"`

	// Other filters - will use Solr
	MinSeats        int     `json:"min_seats,omitempty"`
	MaxPrice        float64 `json:"max_price,omitempty"`
	PetsAllowed     *bool   `json:"pets_allowed,omitempty"`
	SmokingAllowed  *bool   `json:"smoking_allowed,omitempty"`
	MusicAllowed    *bool   `json:"music_allowed,omitempty"`
	MinDriverRating float64 `json:"min_driver_rating,omitempty"`

	// Full-text search
	SearchText string `json:"search_text,omitempty"`

	// Sorting and pagination
	SortBy    string `json:"sort_by,omitempty"` // popularity, price_asc, price_desc, date_asc, date_desc
	SortOrder string `json:"sort_order,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// SearchResponse contains the search results with pagination info
type SearchResponse struct {
	Trips      []*SearchTrip `json:"trips"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// Hash generates a deterministic hash for the query (for caching)
// Same query parameters = same hash
func (q *SearchQuery) Hash() string {
	// Create a normalized copy with only search-relevant fields
	// Note: We exclude Address from Location as it doesn't affect search results
	normalized := struct {
		OriginCity        string
		OriginProvince    string
		OriginLat         float64
		OriginLng         float64
		DestinationCity   string
		DestinationProv   string
		DestinationLat    float64
		DestinationLng    float64
		OriginRadius      int
		DestinationRadius int
		DepartureDate     string
		MinSeats          int
		MaxPrice          float64
		PetsAllowed       *bool
		SmokingAllowed    *bool
		MusicAllowed      *bool
		MinDriverRating   float64
		SearchText        string
		SortBy            string
		SortOrder         string
		Page              int
		Limit             int
	}{
		OriginRadius:      q.OriginRadius,
		DestinationRadius: q.DestinationRadius,
		MinSeats:          q.MinSeats,
		MaxPrice:          q.MaxPrice,
		PetsAllowed:       q.PetsAllowed,
		SmokingAllowed:    q.SmokingAllowed,
		MusicAllowed:      q.MusicAllowed,
		MinDriverRating:   q.MinDriverRating,
		SearchText:        q.SearchText,
		SortBy:            q.SortBy,
		SortOrder:         q.SortOrder,
		Page:              q.Page,
		Limit:             q.Limit,
	}

	// Extract Origin fields if present
	if q.Origin != nil {
		normalized.OriginCity = q.Origin.City
		normalized.OriginProvince = q.Origin.Province
		if len(q.Origin.Coordinates.Coordinates) == 2 {
			normalized.OriginLat = q.Origin.Coordinates.Lat()
			normalized.OriginLng = q.Origin.Coordinates.Lng()
		}
	}

	// Extract Destination fields if present
	if q.Destination != nil {
		normalized.DestinationCity = q.Destination.City
		normalized.DestinationProv = q.Destination.Province
		if len(q.Destination.Coordinates.Coordinates) == 2 {
			normalized.DestinationLat = q.Destination.Coordinates.Lat()
			normalized.DestinationLng = q.Destination.Coordinates.Lng()
		}
	}

	// Format dates consistently (empty string if zero time)
	if q.DepartureDate != nil {
		normalized.DepartureDate = q.DepartureDate.UTC().Format(time.RFC3339)
	}

	// Convert to JSON for consistent representation
	data, _ := json.Marshal(normalized)

	// Generate SHA-256 hash
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// IsGeospatial returns true if this is a geospatial query with radius
// Note: User can provide coordinates without radius (for exact city match)
func (q *SearchQuery) IsGeospatial() bool {
	hasOriginGeo := q.Origin != nil &&
		len(q.Origin.Coordinates.Coordinates) == 2 &&
		q.OriginRadius > 0

	hasDestGeo := q.Destination != nil &&
		len(q.Destination.Coordinates.Coordinates) == 2 &&
		q.DestinationRadius > 0

	return hasOriginGeo || hasDestGeo
}

// Validate checks if the query parameters are valid
// Allows searches without origin/destination to show all available trips
func (q *SearchQuery) Validate() error {

	// Allow searches without origin/destination (show all trips)
	// Just validate that if provided, they have required data

	// Validate Origin
	if q.Origin != nil {
		hasCity := q.Origin.City != ""
		hasCoords := len(q.Origin.Coordinates.Coordinates) == 2

		// Must have at least city or coordinates
		if !hasCity && !hasCoords {
			return fmt.Errorf("origin must have city or coordinates")
		}

		// If has coordinates, validate them
		if hasCoords {
			originLat, originLng := q.Origin.Coordinates.Lat(), q.Origin.Coordinates.Lng()
			if originLat < -90 || originLat > 90 {
				return fmt.Errorf("origin latitude must be between -90 and 90")
			}
			if originLng < -180 || originLng > 180 {
				return fmt.Errorf("origin longitude must be between -180 and 180")
			}
			// Radius is optional now - user can provide coordinates without radius
		} else if q.OriginRadius > 0 {
			// Error: radius without coordinates
			return fmt.Errorf("origin coordinates required when radius specified")
		}
	}

	// Validate Destination
	if q.Destination != nil {
		hasCity := q.Destination.City != ""
		hasCoords := len(q.Destination.Coordinates.Coordinates) == 2

		// Must have at least city or coordinates
		if !hasCity && !hasCoords {
			return fmt.Errorf("destination must have city or coordinates")
		}

		// If has coordinates, validate them
		if hasCoords {
			destinationLat, destinationLng := q.Destination.Coordinates.Lat(), q.Destination.Coordinates.Lng()
			if destinationLat < -90 || destinationLat > 90 {
				return fmt.Errorf("destination latitude must be between -90 and 90")
			}
			if destinationLng < -180 || destinationLng > 180 {
				return fmt.Errorf("destination longitude must be between -180 and 180")
			}
			// Radius is optional now - user can provide coordinates without radius
		} else if q.DestinationRadius > 0 {
			// Error: radius without coordinates
			return fmt.Errorf("destination coordinates required when radius specified")
		}
	}

	// Validate numeric filters
	if q.MinSeats < 0 {
		return fmt.Errorf("min_seats cannot be negative")
	}
	if q.MaxPrice < 0 {
		return fmt.Errorf("max_price cannot be negative")
	}
	if q.MinDriverRating < 0 || q.MinDriverRating > 5 {
		return fmt.Errorf("min_driver_rating must be between 0 and 5")
	}
	// Validate sort_by parameter
	// Supports both new flexible format (price, departure_time, rating, popularity)
	// and old shortcuts for backward compatibility (earliest, cheapest, best_rated)
	validSorts := map[string]bool{
		// New flexible format
		"price":          true,
		"departure_time": true,
		"rating":         true,
		"popularity":     true,
		// Backward compatibility shortcuts
		"earliest":    true,
		"cheapest":    true,
		"best_rated":  true,
	}
	if q.SortBy != "" && !validSorts[q.SortBy] {
		return fmt.Errorf("invalid sort_by")
	}

	// Validate sort_order parameter
	validSortOrders := map[string]bool{"asc": true, "desc": true}
	if q.SortOrder != "" && !validSortOrders[q.SortOrder] {
		return fmt.Errorf("invalid sort_order: must be 'asc' or 'desc'")
	}

	// Validate pagination
	if q.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}
	if q.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if q.Limit > 100 {
		return fmt.Errorf("limit cannot exceed 100")
	}

	return nil
}

// SetDefaults sets default values for pagination if not specified
func (q *SearchQuery) SetDefaults() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 20
	}
	if q.SortBy == "" {
		q.SortBy = "popularity" // Default to most popular
	}
}
