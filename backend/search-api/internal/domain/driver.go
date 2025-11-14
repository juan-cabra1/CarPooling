package domain

// Driver represents denormalized driver information embedded in SearchTrip
// This data is fetched from users-api during event processing and stored here for performance
type Driver struct {
	ID         int64   `json:"id" bson:"id"`
	Name       string  `json:"name" bson:"name"`
	Email      string  `json:"email" bson:"email"`
	PhotoURL   string  `json:"photo_url,omitempty" bson:"photo_url,omitempty"`
	Rating     float64 `json:"rating" bson:"rating"`
	TotalTrips int     `json:"total_trips" bson:"total_trips"`
}
