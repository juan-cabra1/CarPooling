package domain

// Location represents a geographical location with city, province, address and coordinates
type Location struct {
	City        string      `json:"city" bson:"city" binding:"required"`
	Province    string      `json:"province" bson:"province" binding:"required"`
	Address     string      `json:"address" bson:"address" binding:"required"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates" binding:"required"`
}

// Coordinates represents geographical coordinates
// Note: Domain uses {lat, lng} format
// MongoDB stores as GeoJSON [lng, lat] array - conversion handled in repository layer
type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat" binding:"required"`
	Lng float64 `json:"lng" bson:"lng" binding:"required"`
}
