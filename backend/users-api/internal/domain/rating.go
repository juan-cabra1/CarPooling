package domain

import "time"

// RatingDTO representa una calificación en el dominio de negocio
type RatingDTO struct {
	ID          int64     `json:"id"`
	RaterID     int64     `json:"rater_id"`
	RatedUserID int64     `json:"rated_user_id"`
	TripID      string    `json:"trip_id"`
	RoleRated   string    `json:"role_rated"`
	Score       int       `json:"score"`
	Comment     string    `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateRatingRequest representa los datos necesarios para crear una calificación
type CreateRatingRequest struct {
	RaterID     int64  `json:"rater_id" binding:"required"`
	RatedUserID int64  `json:"rated_user_id" binding:"required"`
	TripID      string `json:"trip_id" binding:"required"`
	RoleRated   string `json:"role_rated" binding:"required,oneof=conductor pasajero"`
	Score       int    `json:"score" binding:"required,min=1,max=5"`
	Comment     string `json:"comment"`
}
