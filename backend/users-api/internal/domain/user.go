package domain

import "time"

// User representa un usuario en el dominio de negocio
type User struct {
	ID                  uint64    `json:"id"`
	Email               string    `json:"email"`
	EmailVerified       bool      `json:"email_verified"`
	Name                string    `json:"name"`
	Lastname            string    `json:"lastname"`
	Role                string    `json:"role"`
	Phone               string    `json:"phone"`
	Street              string    `json:"street"`
	Number              int       `json:"number"`
	PhotoUrl            string    `json:"photo_url,omitempty"`
	Sex                 string    `json:"sex"`
	AvgDriverRating     float32   `json:"avg_driver_rating"`
	AvgPassengerRating  float32   `json:"avg_passenger_rating"`
	TotalTripsPassenger uint64    `json:"total_trips_passenger"`
	TotalTripsDriver    uint64    `json:"total_trips_driver"`
	Birthdate           time.Time `json:"birthdate"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// UserCreateRequest representa los datos necesarios para crear un usuario
type UserCreateRequest struct {
	Email     string    `json:"email" validate:"required,email"`
	Name      string    `json:"name" validate:"required"`
	Lastname  string    `json:"lastname" validate:"required"`
	Password  string    `json:"password" validate:"required,min=8"`
	Phone     string    `json:"phone" validate:"required"`
	Street    string    `json:"street" validate:"required"`
	Number    int       `json:"number" validate:"required"`
	PhotoUrl  string    `json:"photo_url,omitempty"`
	Sex       string    `json:"sex" validate:"required,oneof=hombre mujer otro"`
	Birthdate time.Time `json:"birthdate" validate:"required"`
}

// UserUpdateRequest representa los datos que se pueden actualizar de un usuario
type UserUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Lastname *string `json:"lastname,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Street   *string `json:"street,omitempty"`
	Number   *int    `json:"number,omitempty"`
	PhotoUrl *string `json:"photo_url,omitempty"`
}

// UserLoginRequest representa las credenciales de login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserLoginResponse representa la respuesta al login
type UserLoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
