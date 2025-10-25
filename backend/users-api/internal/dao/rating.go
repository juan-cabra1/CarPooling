package dao

import "time"

// RatingDAO representa la estructura de datos para la tabla ratings en MySQL
type RatingDAO struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;column:id"`
	RaterID     int64     `gorm:"not null;index;column:rater_id"`
	RatedUserID int64     `gorm:"not null;index;column:rated_user_id"`
	TripID      string    `gorm:"type:varchar(24);not null;index;column:trip_id"`
	RoleRated   string    `gorm:"type:enum('conductor','pasajero');not null;column:role_rated"`
	Score       int       `gorm:"not null;column:score"`
	Comment     string    `gorm:"type:text;column:comment"`
	CreatedAt   time.Time `gorm:"autoCreateTime;column:created_at"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (RatingDAO) TableName() string {
	return "ratings"
}
