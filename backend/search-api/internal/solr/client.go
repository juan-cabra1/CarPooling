package solr

import (
	"fmt"
	"search-api/internal/domain"
	"time"

	"github.com/rs/zerolog/log"
	solr "github.com/rtt/Go-Solr"
)

// Client provides methods to interact with Apache Solr
type Client struct {
	conn *solr.Connection
	core string
}

// NewClient creates a new Solr client
func NewClient(baseURL, core string) (*Client, error) {
	// Parse baseURL to extract host and port
	// baseURL format: http://localhost:8983 or https://host:port
	var host string
	var port int

	// Simple parsing - assuming format http://host:port or https://host:port
	if len(baseURL) > 7 && baseURL[:7] == "http://" {
		baseURL = baseURL[7:]
	} else if len(baseURL) > 8 && baseURL[:8] == "https://" {
		baseURL = baseURL[8:]
	}

	// Default port for Solr
	host = baseURL
	port = 8983

	// Try to parse host:port
	if idx := len(baseURL) - 1; idx > 0 {
		for i := len(baseURL) - 1; i >= 0; i-- {
			if baseURL[i] == ':' {
				host = baseURL[:i]
				fmt.Sscanf(baseURL[i+1:], "%d", &port)
				break
			}
		}
	}

	conn, err := solr.Init(host, port, core)
	if err != nil {
		return nil, fmt.Errorf("failed to create Solr connection: %w", err)
	}

	log.Info().
		Str("solr_url", baseURL).
		Str("core", core).
		Msg("Solr client initialized successfully")

	return &Client{
		conn: conn,
		core: core,
	}, nil
}

// Index adds or updates a trip document in Solr
func (c *Client) Index(trip *domain.SearchTrip) error {
	if trip == nil {
		return fmt.Errorf("trip cannot be nil")
	}

	// Map SearchTrip to Solr document
	doc := MapTripToSolrDocument(trip)

	// Prepare update command
	updateDoc := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": doc,
		},
	}

	// Add document to Solr
	if _, err := c.conn.Update(updateDoc, true); err != nil {
		log.Error().
			Err(err).
			Str("trip_id", trip.TripID).
			Msg("Failed to index trip in Solr")
		return fmt.Errorf("failed to index trip: %w", err)
	}

	log.Debug().
		Str("trip_id", trip.TripID).
		Str("core", c.core).
		Msg("Trip indexed successfully in Solr")

	return nil
}

// Delete removes a trip document from Solr
func (c *Client) Delete(tripID string) error {
	if tripID == "" {
		return fmt.Errorf("tripID cannot be empty")
	}

	// Delete by ID
	deleteDoc := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": tripID,
		},
	}

	// Execute delete
	resp, err := c.conn.Update(deleteDoc, true)
	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to delete trip from Solr")
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	log.Debug().
		Str("trip_id", tripID).
		Interface("response", resp).
		Msg("Trip deleted successfully from Solr")

	return nil
}

// Search performs a search query in Solr
func (c *Client) Search(query *SearchQuery) (*SearchResponse, error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}

	// Build Solr query
	solrQuery := BuildSolrQuery(query)

	// Execute search
	resp, err := c.conn.Select(solrQuery)
	if err != nil {
		log.Error().
			Err(err).
			Interface("query", query).
			Msg("Failed to execute Solr search")
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}

	// Parse response
	searchResp := &SearchResponse{
		NumFound: resp.Results.NumFound,
		Start:    resp.Results.Start,
		Docs:     make([]map[string]interface{}, 0),
	}

	// Extract documents from Collection
	for _, doc := range resp.Results.Collection {
		searchResp.Docs = append(searchResp.Docs, doc.Fields)
	}

	// Extract facets if available
	if resp.Results.Facets != nil && len(resp.Results.Facets) > 0 {
		searchResp.Facets = resp.Results.Facets
	}

	log.Debug().
		Int("num_found", searchResp.NumFound).
		Int("returned", len(searchResp.Docs)).
		Msg("Solr search completed successfully")

	return searchResp, nil
}

// Ping checks if Solr is reachable
func (c *Client) Ping() error {
	// Simple query to test connection
	query := &solr.Query{
		Params: solr.URLParamMap{
			"q":    []string{"*:*"},
			"rows": []string{"0"},
		},
	}

	_, err := c.conn.Select(query)
	if err != nil {
		log.Warn().
			Err(err).
			Str("core", c.core).
			Msg("Solr ping failed")
		return fmt.Errorf("solr ping failed: %w", err)
	}

	log.Debug().
		Str("core", c.core).
		Msg("Solr ping successful")

	return nil
}

// GetConnection returns the underlying Solr connection (for health checks)
func (c *Client) GetConnection() *solr.Connection {
	return c.conn
}

// SearchQuery represents a search query with filters
type SearchQuery struct {
	// Full-text search
	Query string

	// Filters
	OriginCity      string
	DestinationCity string
	MinSeats        int
	MaxPrice        float64
	MinPrice        float64
	Status          string

	// Preferences
	PetsAllowed    *bool
	SmokingAllowed *bool
	MusicAllowed   *bool

	// Date range
	DepartureFrom time.Time
	DepartureTo   time.Time

	// Pagination
	Page  int
	Limit int

	// Sorting
	SortBy    string // e.g., "price_per_seat", "departure_datetime"
	SortOrder string // "asc" or "desc"

	// Facets
	EnableFacets bool
	FacetFields  []string
}

// SearchResponse represents the response from a Solr search
type SearchResponse struct {
	NumFound int                      `json:"numFound"`
	Start    int                      `json:"start"`
	Docs     []map[string]interface{} `json:"docs"`
	Facets   interface{}              `json:"facets,omitempty"`
}
