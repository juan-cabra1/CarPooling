package domain

// Location represents a geographical location with city, province, address and coordinates
type Location struct {
	City        string       `json:"city" bson:"city" binding:"required"`
	Province    string       `json:"province" bson:"province" binding:"required"`
	Address     string       `json:"address" bson:"address" binding:"required"`
	Coordinates GeoJSONPoint `json:"coordinates" bson:"coordinates" binding:"required"`
}

// GeoJSONPoint represents geographical coordinates in GeoJSON format
// This is the standard format required by MongoDB's 2dsphere indexes
// Reference: https://www.mongodb.com/docs/manual/reference/geojson/
type GeoJSONPoint struct {
	Type        string    `json:"type" bson:"type"`                   // Must be "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`     // [longitude, latitude] - IMPORTANT: lng first!
}

// NewGeoJSONPoint creates a GeoJSON Point from latitude and longitude
// Parameters: lat (latitude), lng (longitude)
// Returns GeoJSON Point with coordinates in [lng, lat] format as required by MongoDB
func NewGeoJSONPoint(lat, lng float64) GeoJSONPoint {
	return GeoJSONPoint{
		Type:        "Point",
		Coordinates: []float64{lng, lat}, // MongoDB GeoJSON format: [longitude, latitude]
	}
}

// Lat returns the latitude from GeoJSON coordinates
func (g GeoJSONPoint) Lat() float64 {
	if len(g.Coordinates) >= 2 {
		return g.Coordinates[1] // Latitude is at index 1
	}
	return 0
}

// Lng returns the longitude from GeoJSON coordinates
func (g GeoJSONPoint) Lng() float64 {
	if len(g.Coordinates) >= 1 {
		return g.Coordinates[0] // Longitude is at index 0
	}
	return 0
}
