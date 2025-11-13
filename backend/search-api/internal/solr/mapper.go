package solr

import (
	"fmt"
	"search-api/internal/domain"
	"strings"
	"time"
)

// MapTripToSolrDocument converts a SearchTrip domain model to a Solr document
// NOTE: Geospatial fields are NOT included - MongoDB handles geospatial queries
func MapTripToSolrDocument(trip *domain.SearchTrip) map[string]interface{} {
	doc := make(map[string]interface{})

	// Primary identifier
	doc["id"] = trip.TripID

	// Driver information
	doc["driver_id"] = trip.DriverID
	doc["driver_name"] = trip.Driver.Name
	doc["driver_rating"] = trip.Driver.Rating
	doc["driver_total_trips"] = trip.Driver.TotalTrips

	// Location information (ONLY city and province - NO coordinates)
	doc["origin_city"] = trip.Origin.City
	doc["origin_province"] = trip.Origin.Province
	doc["destination_city"] = trip.Destination.City
	doc["destination_province"] = trip.Destination.Province

	// Trip timing (Solr date format: ISO 8601)
	doc["departure_datetime"] = formatSolrDate(trip.DepartureDatetime)
	doc["estimated_arrival_datetime"] = formatSolrDate(trip.EstimatedArrivalDatetime)

	// Pricing and availability
	doc["price_per_seat"] = trip.PricePerSeat
	doc["total_seats"] = trip.TotalSeats
	doc["available_seats"] = trip.AvailableSeats

	// Vehicle information
	doc["car_brand"] = trip.Car.Brand
	doc["car_model"] = trip.Car.Model
	doc["car_year"] = trip.Car.Year
	doc["car_color"] = trip.Car.Color

	// Preferences (boolean fields)
	doc["pets_allowed"] = trip.Preferences.PetsAllowed
	doc["smoking_allowed"] = trip.Preferences.SmokingAllowed
	doc["music_allowed"] = trip.Preferences.MusicAllowed

	// Trip details
	doc["status"] = trip.Status
	doc["description"] = trip.Description

	// Search-specific fields
	doc["search_text"] = buildSearchText(trip)
	doc["popularity_score"] = trip.PopularityScore

	// Timestamps
	doc["created_at"] = formatSolrDate(trip.CreatedAt)
	doc["updated_at"] = formatSolrDate(trip.UpdatedAt)

	return doc
}

// buildSearchText concatenates all relevant text fields for full-text search
func buildSearchText(trip *domain.SearchTrip) string {
	parts := []string{
		trip.Origin.City,
		trip.Origin.Province,
		trip.Origin.Address,
		trip.Destination.City,
		trip.Destination.Province,
		trip.Destination.Address,
		trip.Driver.Name,
		trip.Car.Brand,
		trip.Car.Model,
		trip.Car.Color,
		trip.Description,
	}

	// Filter out empty strings
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}

	return strings.Join(filtered, " ")
}

// formatSolrDate converts Go time.Time to Solr date format (ISO 8601)
// Example: 2024-01-15T14:30:00Z
func formatSolrDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// MapSolrDocumentToTrip converts a Solr document back to SearchTrip (for results)
// This is a partial mapping - full trip details should come from MongoDB if needed
func MapSolrDocumentToTrip(doc map[string]interface{}) (*domain.SearchTrip, error) {
	trip := &domain.SearchTrip{}

	// Extract basic fields
	if id, ok := doc["id"].(string); ok {
		trip.TripID = id
	}

	if driverID, ok := doc["driver_id"].(float64); ok {
		trip.DriverID = int64(driverID)
	}

	// Driver information
	trip.Driver = domain.Driver{
		Name: getStringField(doc, "driver_name"),
	}
	if rating, ok := doc["driver_rating"].(float64); ok {
		trip.Driver.Rating = rating
	}
	if totalTrips, ok := doc["driver_total_trips"].(float64); ok {
		trip.Driver.TotalTrips = int(totalTrips)
	}

	// Location information
	trip.Origin = domain.Location{
		City:     getStringField(doc, "origin_city"),
		Province: getStringField(doc, "origin_province"),
	}
	trip.Destination = domain.Location{
		City:     getStringField(doc, "destination_city"),
		Province: getStringField(doc, "destination_province"),
	}

	// Timing
	if depTime, err := parseSolrDate(getStringField(doc, "departure_datetime")); err == nil {
		trip.DepartureDatetime = depTime
	}
	if arrTime, err := parseSolrDate(getStringField(doc, "estimated_arrival_datetime")); err == nil {
		trip.EstimatedArrivalDatetime = arrTime
	}

	// Pricing and availability
	if price, ok := doc["price_per_seat"].(float64); ok {
		trip.PricePerSeat = price
	}
	if total, ok := doc["total_seats"].(float64); ok {
		trip.TotalSeats = int(total)
	}
	if avail, ok := doc["available_seats"].(float64); ok {
		trip.AvailableSeats = int(avail)
	}

	// Vehicle
	trip.Car = domain.Car{
		Brand: getStringField(doc, "car_brand"),
		Model: getStringField(doc, "car_model"),
		Color: getStringField(doc, "car_color"),
	}
	if year, ok := doc["car_year"].(float64); ok {
		trip.Car.Year = int(year)
	}

	// Preferences
	trip.Preferences = domain.Preferences{
		PetsAllowed:    getBoolField(doc, "pets_allowed"),
		SmokingAllowed: getBoolField(doc, "smoking_allowed"),
		MusicAllowed:   getBoolField(doc, "music_allowed"),
	}

	// Status and description
	trip.Status = getStringField(doc, "status")
	trip.Description = getStringField(doc, "description")

	// Search fields
	if score, ok := doc["popularity_score"].(float64); ok {
		trip.PopularityScore = score
	}

	return trip, nil
}

// Helper functions for type-safe field extraction

func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	return ""
}

func getBoolField(doc map[string]interface{}, field string) bool {
	if val, ok := doc[field].(bool); ok {
		return val
	}
	return false
}

func parseSolrDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}
	return time.Parse(time.RFC3339, dateStr)
}
