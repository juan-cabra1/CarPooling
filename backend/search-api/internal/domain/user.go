package domain

import "time"

// User represents the external User DTO from users-api
// This is the structure returned by users-api GET /users/:id endpoint
type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth"`

	// Verification status
	IsVerified      bool `json:"is_verified"`
	IsEmailVerified bool `json:"is_email_verified"`
	IsPhoneVerified bool `json:"is_phone_verified"`

	// Trip statistics
	TotalTripsAsDriver    int `json:"total_trips_as_driver"`
	TotalTripsAsPassenger int `json:"total_trips_as_passenger"`

	// Rating information
	AverageRatingAsDriver    float64 `json:"average_rating_as_driver"`
	TotalRatingsAsDriver     int     `json:"total_ratings_as_driver"`
	AverageRatingAsPassenger float64 `json:"average_rating_as_passenger"`
	TotalRatingsAsPassenger  int     `json:"total_ratings_as_passenger"`

	// Preferences
	PreferredLanguage string `json:"preferred_language,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToDriver converts User DTO to Driver (embedded in SearchTrip)
// Extracts only the fields needed for search denormalization
func (u *User) ToDriver() Driver {
	return Driver{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		PhotoURL:   u.PhotoURL,
		Rating:     u.AverageRatingAsDriver,
		TotalTrips: u.TotalTripsAsDriver,
	}
}
