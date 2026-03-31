package service

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles JWT token validation
type AuthService interface {
	// ValidateToken validates a JWT token and returns the parsed token
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// authService implements AuthService
type authService struct {
	jwtSecret string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

// ValidateToken validates a JWT token and returns the parsed token
//
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - *jwt.Token: The parsed and validated token
//   - error: Error if validation fails
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
