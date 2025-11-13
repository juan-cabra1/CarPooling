package solr

import (
	"fmt"
	"strings"

	solr "github.com/rtt/Go-Solr"
)

// addParam is a helper function to add a parameter to URLParamMap
func addParam(params solr.URLParamMap, key, value string) {
	params[key] = append(params[key], value)
}

// BuildSolrQuery constructs a Solr query from SearchQuery parameters
func BuildSolrQuery(sq *SearchQuery) *solr.Query {
	query := &solr.Query{
		Params: solr.URLParamMap{},
	}

	// Main query (q parameter)
	mainQuery := buildMainQuery(sq)
	addParam(query.Params, "q", mainQuery)

	// Filter queries (fq parameters) - more efficient for filtering
	filterQueries := buildFilterQueries(sq)
	for _, fq := range filterQueries {
		addParam(query.Params, "fq", fq)
	}

	// Pagination
	start := (sq.Page - 1) * sq.Limit
	if start < 0 {
		start = 0
	}
	addParam(query.Params, "start", fmt.Sprintf("%d", start))
	addParam(query.Params, "rows", fmt.Sprintf("%d", sq.Limit))

	// Sorting
	if sq.SortBy != "" {
		sortOrder := "asc"
		if sq.SortOrder == "desc" {
			sortOrder = "desc"
		}
		addParam(query.Params, "sort", fmt.Sprintf("%s %s", sq.SortBy, sortOrder))
	}

	// Facets
	if sq.EnableFacets {
		addParam(query.Params, "facet", "true")
		if len(sq.FacetFields) > 0 {
			for _, field := range sq.FacetFields {
				addParam(query.Params, "facet.field", field)
			}
		} else {
			// Default facet fields
			addParam(query.Params, "facet.field", "origin_city")
			addParam(query.Params, "facet.field", "destination_city")
			addParam(query.Params, "facet.field", "status")
		}
		addParam(query.Params, "facet.mincount", "1")
	}

	// Response format
	addParam(query.Params, "wt", "json")

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
	addParam(query.Params, "q", queryStr)
	addParam(query.Params, "rows", fmt.Sprintf("%d", limit))
	addParam(query.Params, "wt", "json")

	// Return only the relevant field
	addParam(query.Params, "fl", field)

	// Group by field to get unique values
	addParam(query.Params, "group", "true")
	addParam(query.Params, "group.field", field)
	addParam(query.Params, "group.limit", "1")

	return query
}

// BuildFacetQuery constructs a query specifically for getting facet counts
func BuildFacetQuery(facetFields []string) *solr.Query {
	query := &solr.Query{
		Params: solr.URLParamMap{},
	}

	addParam(query.Params, "q", "*:*")
	addParam(query.Params, "rows", "0") // We only want facets, not documents
	addParam(query.Params, "facet", "true")
	addParam(query.Params, "facet.mincount", "1")

	for _, field := range facetFields {
		addParam(query.Params, "facet.field", field)
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
	addParam(query.Params, "json.facet", jsonFacets)

	addParam(query.Params, "wt", "json")

	return query
}
