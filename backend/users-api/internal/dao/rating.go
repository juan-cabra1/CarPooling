package dao

import "time"

type Rating struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	RaterID     uint64 `gorm:"not null;index"` // Usuario que califica
	RatedUserID uint64 `gorm:"not null;index"` // Usuario calificado
	TripID      uint64 `gorm:"not null;index"` // Viaje relacionado
	RoleRated   string `gorm:"type:enum('conductor','pasajero');not null"`
	Score       uint   `gorm:"check:score >= 1 AND score <= 5;not null"`
	Comment     string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
