package service

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService define las operaciones de autenticación
type AuthService interface {
	// ValidateToken valida un token JWT y retorna los claims
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type authService struct {
	jwtSecret string
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

// ValidateToken valida un token JWT
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parsear el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
