package dao

import "time"

type User struct {
	ID                  uint64 `gorm:"primaryKey;autoIncrement"`
	Email               string `gorm:"unique;not null"`
	EmailVerified       bool   `gorm:"default:false;not null"`
	Name                string `gorm:"not null"`
	Lastname            string `gorm:"not null"`
	PasswordHash        string `gorm:"not null"`
	Role                string `gorm:"enum('admin', 'user');not null"`
	Phone               string `gorm:"not null"`
	Streeat             string `gorm:"not null"`
	Number              int    `gorm:"not null"`
	PhotoUrl            string
	Sex                 string `gorm:"enum('hombre', 'mujer', 'otro');not null"`
	AvgDriverRating     float32
	AvgPassengerRating  float32
	TotalTripsPassenger uint64
	TotalTripsDriver    uint64
	Birthdate           time.Time `gorm:"not null"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}
