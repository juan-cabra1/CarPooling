package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// SearchQuery contains all filter parameters for searching trips
type SearchQuery struct {
	// City-based search
	OriginCity      string `json:"origin_city,omitempty"`
	DestinationCity string `json:"destination_city,omitempty"`

	// Geospatial params - will use MongoDB 2dsphere
	OriginLat      float64 `json:"origin_lat,omitempty"`
	OriginLng      float64 `json:"origin_lng,omitempty"`
	OriginRadius   int     `json:"origin_radius,omitempty"` // in kilometers
	DestinationLat float64 `json:"destination_lat,omitempty"`
	DestinationLng float64 `json:"destination_lng,omitempty"`
	DestinationRadius int  `json:"destination_radius,omitempty"` // in kilometers

	// Date range filters
	DateFrom time.Time `json:"date_from,omitempty"`
	DateTo   time.Time `json:"date_to,omitempty"`

	// Other filters - will use Solr
	MinSeats         int     `json:"min_seats,omitempty"`
	MaxPrice         float64 `json:"max_price,omitempty"`
	PetsAllowed      *bool   `json:"pets_allowed,omitempty"`
	SmokingAllowed   *bool   `json:"smoking_allowed,omitempty"`
	MusicAllowed     *bool   `json:"music_allowed,omitempty"`
	MinDriverRating  float64 `json:"min_driver_rating,omitempty"`

	// Full-text search
	SearchText string `json:"search_text,omitempty"`

	// Sorting and pagination
	SortBy string `json:"sort_by,omitempty"` // popularity, price_asc, price_desc, date_asc, date_desc
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
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
	// Create a normalized copy for consistent hashing
	normalized := struct {
		OriginCity        string
		DestinationCity   string
		OriginLat         float64
		OriginLng         float64
		OriginRadius      int
		DestinationLat    float64
		DestinationLng    float64
		DestinationRadius int
		DateFrom          string
		DateTo            string
		MinSeats          int
		MaxPrice          float64
		PetsAllowed       *bool
		SmokingAllowed    *bool
		MusicAllowed      *bool
		MinDriverRating   float64
		SearchText        string
		SortBy            string
		Page              int
		Limit             int
	}{
		OriginCity:        q.OriginCity,
		DestinationCity:   q.DestinationCity,
		OriginLat:         q.OriginLat,
		OriginLng:         q.OriginLng,
		OriginRadius:      q.OriginRadius,
		DestinationLat:    q.DestinationLat,
		DestinationLng:    q.DestinationLng,
		DestinationRadius: q.DestinationRadius,
		MinSeats:          q.MinSeats,
		MaxPrice:          q.MaxPrice,
		PetsAllowed:       q.PetsAllowed,
		SmokingAllowed:    q.SmokingAllowed,
		MusicAllowed:      q.MusicAllowed,
		MinDriverRating:   q.MinDriverRating,
		SearchText:        q.SearchText,
		SortBy:            q.SortBy,
		Page:              q.Page,
		Limit:             q.Limit,
	}

	// Format dates consistently (empty string if zero time)
	if !q.DateFrom.IsZero() {
		normalized.DateFrom = q.DateFrom.UTC().Format(time.RFC3339)
	}
	if !q.DateTo.IsZero() {
		normalized.DateTo = q.DateTo.UTC().Format(time.RFC3339)
	}

	// Convert to JSON for consistent representation
	data, _ := json.Marshal(normalized)

	// Generate SHA-256 hash
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// IsGeospatial returns true if this is a geospatial query
func (q *SearchQuery) IsGeospatial() bool {
	return (q.OriginLat != 0 && q.OriginLng != 0) ||
		(q.DestinationLat != 0 && q.DestinationLng != 0)
}

// Validate checks if the query parameters are valid
func (q *SearchQuery) Validate() error {
	// Validate geospatial coordinates
	if q.OriginLat != 0 || q.OriginLng != 0 {
		if q.OriginLat < -90 || q.OriginLat > 90 {
			return fmt.Errorf("origin latitude must be between -90 and 90")
		}
		if q.OriginLng < -180 || q.OriginLng > 180 {
			return fmt.Errorf("origin longitude must be between -180 and 180")
		}
		if q.OriginRadius <= 0 {
			return fmt.Errorf("origin radius must be positive")
		}
	}

	if q.DestinationLat != 0 || q.DestinationLng != 0 {
		if q.DestinationLat < -90 || q.DestinationLat > 90 {
			return fmt.Errorf("destination latitude must be between -90 and 90")
		}
		if q.DestinationLng < -180 || q.DestinationLng > 180 {
			return fmt.Errorf("destination longitude must be between -180 and 180")
		}
		if q.DestinationRadius <= 0 {
			return fmt.Errorf("destination radius must be positive")
		}
	}

	// Validate date range
	if !q.DateFrom.IsZero() && !q.DateTo.IsZero() {
		if q.DateFrom.After(q.DateTo) {
			return fmt.Errorf("date_from must be before date_to")
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
