package domain

import "time"

// UserDTO representa un usuario en el dominio de negocio
type UserDTO struct {
	ID                  int64     `json:"id"`
	Email               string    `json:"email"`
	EmailVerified       bool      `json:"email_verified"`
	Name                string    `json:"name"`
	Lastname            string    `json:"lastname"`
	Role                string    `json:"role"`
	Phone               string    `json:"phone"`
	Street              string    `json:"street"`
	Number              int       `json:"number"`
	PhotoURL            string    `json:"photo_url,omitempty"`
	Sex                 string    `json:"sex"`
	AvgDriverRating     float64   `json:"avg_driver_rating"`
	AvgPassengerRating  float64   `json:"avg_passenger_rating"`
	TotalTripsPassenger int       `json:"total_trips_passenger"`
	TotalTripsDriver    int       `json:"total_trips_driver"`
	Birthdate           time.Time `json:"birthdate"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// CreateUserRequest representa los datos necesarios para crear un usuario
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Name      string `json:"name" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Street    string `json:"street" binding:"required"`
	Number    int    `json:"number" binding:"required"`
	PhotoURL  string `json:"photo_url"`
	Sex       string `json:"sex" binding:"required,oneof=hombre mujer otro"`
	Birthdate string `json:"birthdate" binding:"required"` // Format: YYYY-MM-DD
}

// UpdateUserRequest representa los datos que se pueden actualizar de un usuario
type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Lastname *string `json:"lastname"`
	Phone    *string `json:"phone"`
	Street   *string `json:"street"`
	Number   *int    `json:"number"`
	PhotoURL *string `json:"photo_url"`
}

// LoginRequest representa las credenciales de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse representa la respuesta al login
type LoginResponse struct {
	Token string   `json:"token"`
	User  *UserDTO `json:"user"`
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest representa la solicitud para restablecer contraseña
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResendVerificationRequest representa la solicitud para reenviar email de verificación
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
