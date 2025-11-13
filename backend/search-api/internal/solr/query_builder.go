package solr

import (
	"fmt"
	"strings"

	solr "github.com/rtt/Go-Solr"
)

// BuildSolrQuery constructs a Solr query from SearchQuery parameters
func BuildSolrQuery(sq *SearchQuery) *solr.Query {
	query := &solr.Query{
		Params: solr.URLParamMap{},
	}

	// Main query (q parameter)
	mainQuery := buildMainQuery(sq)
	query.Params.Add("q", mainQuery)

	// Filter queries (fq parameters) - more efficient for filtering
	filterQueries := buildFilterQueries(sq)
	for _, fq := range filterQueries {
		query.Params.Add("fq", fq)
	}

	// Pagination
	start := (sq.Page - 1) * sq.Limit
	if start < 0 {
		start = 0
	}
	query.Params.Add("start", fmt.Sprintf("%d", start))
	query.Params.Add("rows", fmt.Sprintf("%d", sq.Limit))

	// Sorting
	if sq.SortBy != "" {
		sortOrder := "asc"
		if sq.SortOrder == "desc" {
			sortOrder = "desc"
		}
		query.Params.Add("sort", fmt.Sprintf("%s %s", sq.SortBy, sortOrder))
	}

	// Facets
	if sq.EnableFacets {
		query.Params.Add("facet", "true")
		if len(sq.FacetFields) > 0 {
			for _, field := range sq.FacetFields {
				query.Params.Add("facet.field", field)
			}
		} else {
			// Default facet fields
			query.Params.Add("facet.field", "origin_city")
			query.Params.Add("facet.field", "destination_city")
			query.Params.Add("facet.field", "status")
		}
		query.Params.Add("facet.mincount", "1")
	}

	// Response format
	query.Params.Add("wt", "json")

	return query
}

// buildMainQuery constructs the main query (q parameter)
func buildMainQuery(sq *SearchQuery) string {
	// If there's a search query, use it for full-text search
	if sq.Query != "" {
		// Search in description and search_text fields
		return fmt.Sprintf("(description:%s OR search_text:%s)",
			escapeSolrQuery(sq.Query),
			escapeSolrQuery(sq.Query))
	}

	// Default: match all documents
	return "*:*"
}

// buildFilterQueries constructs filter queries (fq parameters)
// Filter queries are cached by Solr and are more efficient than including in main query
func buildFilterQueries(sq *SearchQuery) []string {
	filters := []string{}

	// City filters
	if sq.OriginCity != "" {
		filters = append(filters, fmt.Sprintf("origin_city:\"%s\"", escapeSolrValue(sq.OriginCity)))
	}
	if sq.DestinationCity != "" {
		filters = append(filters, fmt.Sprintf("destination_city:\"%s\"", escapeSolrValue(sq.DestinationCity)))
	}

	// Availability filter
	if sq.MinSeats > 0 {
		filters = append(filters, fmt.Sprintf("available_seats:[%d TO *]", sq.MinSeats))
	}

	// Price range filter
	if sq.MinPrice > 0 && sq.MaxPrice > 0 {
		filters = append(filters, fmt.Sprintf("price_per_seat:[%f TO %f]", sq.MinPrice, sq.MaxPrice))
	} else if sq.MinPrice > 0 {
		filters = append(filters, fmt.Sprintf("price_per_seat:[%f TO *]", sq.MinPrice))
	} else if sq.MaxPrice > 0 {
		filters = append(filters, fmt.Sprintf("price_per_seat:[* TO %f]", sq.MaxPrice))
	}

	// Status filter
	if sq.Status != "" {
		filters = append(filters, fmt.Sprintf("status:\"%s\"", escapeSolrValue(sq.Status)))
	}

	// Date range filter
	if !sq.DepartureFrom.IsZero() && !sq.DepartureTo.IsZero() {
		fromStr := sq.DepartureFrom.UTC().Format("2006-01-02T15:04:05Z")
		toStr := sq.DepartureTo.UTC().Format("2006-01-02T15:04:05Z")
		filters = append(filters, fmt.Sprintf("departure_datetime:[%s TO %s]", fromStr, toStr))
	} else if !sq.DepartureFrom.IsZero() {
		fromStr := sq.DepartureFrom.UTC().Format("2006-01-02T15:04:05Z")
		filters = append(filters, fmt.Sprintf("departure_datetime:[%s TO *]", fromStr))
	} else if !sq.DepartureTo.IsZero() {
		toStr := sq.DepartureTo.UTC().Format("2006-01-02T15:04:05Z")
		filters = append(filters, fmt.Sprintf("departure_datetime:[* TO %s]", toStr))
	}

	// Preference filters (boolean)
	if sq.PetsAllowed != nil {
		filters = append(filters, fmt.Sprintf("pets_allowed:%t", *sq.PetsAllowed))
	}
	if sq.SmokingAllowed != nil {
		filters = append(filters, fmt.Sprintf("smoking_allowed:%t", *sq.SmokingAllowed))
	}
	if sq.MusicAllowed != nil {
		filters = append(filters, fmt.Sprintf("music_allowed:%t", *sq.MusicAllowed))
	}

	return filters
}

// escapeSolrQuery escapes special characters in a Solr query string
// Solr special characters: + - && || ! ( ) { } [ ] ^ " ~ * ? : \ /
func escapeSolrQuery(s string) string {
	specialChars := []string{
		"\\", "+", "-", "&&", "||", "!", "(", ")", "{", "}",
		"[", "]", "^", "\"", "~", "*", "?", ":", "/",
	}

	result := s
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	return result
}

// escapeSolrValue escapes quotes in field values
func escapeSolrValue(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// BuildAutocompleteQuery constructs a query for autocomplete suggestions
func BuildAutocompleteQuery(prefix string, field string, limit int) *solr.Query {
	query := &solr.Query{
		Params: solr.URLParamMap{},
	}

	// Use wildcard query for prefix matching
	queryStr := fmt.Sprintf("%s:%s*", field, escapeSolrQuery(prefix))
	query.Params.Add("q", queryStr)
	query.Params.Add("rows", fmt.Sprintf("%d", limit))
	query.Params.Add("wt", "json")

	// Return only the relevant field
	query.Params.Add("fl", field)

	// Group by field to get unique values
	query.Params.Add("group", "true")
	query.Params.Add("group.field", field)
	query.Params.Add("group.limit", "1")

	return query
}

// BuildFacetQuery constructs a query specifically for getting facet counts
func BuildFacetQuery(facetFields []string) *solr.Query {
	query := &solr.Query{
		Params: solr.URLParamMap{},
	}

	query.Params.Add("q", "*:*")
	query.Params.Add("rows", "0") // We only want facets, not documents
	query.Params.Add("facet", "true")
	query.Params.Add("facet.mincount", "1")

	for _, field := range facetFields {
		query.Params.Add("facet.field", field)
	}

	// JSON facets for more advanced aggregations
	// Example: price ranges
	jsonFacets := `{
		price_ranges: {
			type: range,
			field: price_per_seat,
			start: 0,
			end: 100000,
			gap: 10000
		},
		popular_routes: {
			type: query,
			q: "*:*",
			facet: {
				routes: {
					type: terms,
					field: origin_city,
					limit: 10,
					facet: {
						destinations: {
							type: terms,
							field: destination_city,
							limit: 10
						}
					}
				}
			}
		}
	}`
	query.Params.Add("json.facet", jsonFacets)

	query.Params.Add("wt", "json")

	return query
}
