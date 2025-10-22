package domain

import "time"

// Rating representa una calificación en el dominio de negocio
type Rating struct {
	ID          uint64    `json:"id"`
	RaterID     uint64    `json:"rater_id"`
	RatedUserID uint64    `json:"rated_user_id"`
	TripID      uint64    `json:"trip_id"`
	RoleRated   string    `json:"role_rated"`
	Score       uint      `json:"score"`
	Comment     string    `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// RatingCreateRequest representa los datos necesarios para crear una calificación
type RatingCreateRequest struct {
	RatedUserID uint64 `json:"rated_user_id" validate:"required"`
	TripID      uint64 `json:"trip_id" validate:"required"`
	RoleRated   string `json:"role_rated" validate:"required,oneof=conductor pasajero"`
	Score       uint   `json:"score" validate:"required,min=1,max=5"`
	Comment     string `json:"comment,omitempty"`
}

// UserRatingSummary representa el resumen de calificaciones de un usuario
type UserRatingSummary struct {
	UserID                uint64  `json:"user_id"`
	AvgDriverRating       float32 `json:"avg_driver_rating"`
	AvgPassengerRating    float32 `json:"avg_passenger_rating"`
	TotalDriverRatings    int64   `json:"total_driver_ratings"`
	TotalPassengerRatings int64   `json:"total_passenger_ratings"`
}
