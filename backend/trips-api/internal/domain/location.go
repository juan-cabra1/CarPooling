package domain

// Location representa una ubicación geográfica con coordenadas
type Location struct {
	City        string      `json:"city" bson:"city" binding:"required"`
	Province    string      `json:"province" bson:"province" binding:"required"`
	Address     string      `json:"address" bson:"address" binding:"required"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates" binding:"required"`
}

// Coordinates representa las coordenadas geográficas (latitud y longitud)
type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat" binding:"required"`
	Lng float64 `json:"lng" bson:"lng" binding:"required"`
}
