package dao

import "time"

// UserDAO representa la estructura de datos para la tabla users en MySQL
type UserDAO struct {
	ID                     int64      `gorm:"primaryKey;autoIncrement;column:id"`
	Email                  string     `gorm:"type:varchar(255);unique;not null;index;column:email"`
	EmailVerified          bool       `gorm:"default:false;not null;column:email_verified"`
	EmailVerificationToken *string    `gorm:"type:varchar(255);column:email_verification_token;index:idx_email_verification_token"`
	PasswordResetToken     *string    `gorm:"type:varchar(255);column:password_reset_token;index:idx_password_reset_token"`
	PasswordResetExpires   *time.Time `gorm:"column:password_reset_expires"`
	Name         string `gorm:"type:varchar(100);not null;column:name"`
	Lastname     string `gorm:"type:varchar(100);not null;column:lastname"`
	PasswordHash string `gorm:"type:varchar(255);not null;column:password_hash"`
	Role         string `gorm:"type:enum('user','admin');default:'user';not null;column:role"`
	Phone        string `gorm:"type:varchar(20);not null;column:phone"`
	Street       string `gorm:"type:varchar(255);not null;column:street"`
	Number       int    `gorm:"not null;column:number"`
	PhotoURL     string `gorm:"type:varchar(255);column:photo_url"`
	Sex          string `gorm:"type:enum('hombre','mujer','otro');not null;column:sex"`
	AvgDriverRating       float64    `gorm:"type:decimal(3,2);default:0.00;column:avg_driver_rating"`
	AvgPassengerRating    float64    `gorm:"type:decimal(3,2);default:0.00;column:avg_passenger_rating"`
	TotalTripsPassenger   int        `gorm:"default:0;column:total_trips_passenger"`
	TotalTripsDriver      int        `gorm:"default:0;column:total_trips_driver"`
	Birthdate             time.Time  `gorm:"not null;column:birthdate"`
	CreatedAt             time.Time  `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime;column:updated_at"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (UserDAO) TableName() string {
	return "users"
}
