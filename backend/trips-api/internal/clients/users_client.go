package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"trips-api/internal/domain"
)

// User representa la información básica de un usuario obtenida desde users-api
type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UsersClient define las operaciones para interactuar con users-api
type UsersClient interface {
	// GetUser obtiene información de un usuario por su ID
	// authToken: JWT token para autenticación en users-api (format: "Bearer {token}")
	// Retorna domain.ErrDriverNotFound si el usuario no existe
	GetUser(ctx context.Context, userID int64, authToken string) (*User, error)
}

type usersHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUsersClient crea una nueva instancia del cliente HTTP para users-api
func NewUsersClient(baseURL string) UsersClient {
	return &usersHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Timeout de 5 segundos para llamadas HTTP
		},
	}
}

// usersAPIResponse representa la estructura de respuesta estándar de users-api
// Formato: {"success": true, "data": {...}}
type usersAPIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// GetUser obtiene información de un usuario desde users-api
//
// Endpoint: GET {base_url}/users/{userID}
// Response: {"success": true, "data": {"id": 123, "name": "John", "email": "john@example.com"}}
//
// Manejo de errores:
// - 404 → domain.ErrDriverNotFound
// - Timeout/Network → error wrapped
// - 5xx → error de servicio externo
//
// Ejemplo de uso:
//
//	user, err := client.GetUser(ctx, 123)
//	if err != nil {
//	    if errors.Is(err, domain.ErrDriverNotFound) {
//	        return fmt.Errorf("driver does not exist")
//	    }
//	    return fmt.Errorf("failed to validate driver: %w", err)
//	}
func (c *usersHTTPClient) GetUser(ctx context.Context, userID int64, authToken string) (*User, error) {
	// Construir URL del endpoint
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	// Crear request con contexto
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Forward Authorization header to users-api
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	// Ejecutar request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Error de red o timeout
		return nil, fmt.Errorf("failed to call users-api: %w", err)
	}
	defer resp.Body.Close()

	// Leer body de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Manejar códigos de estado HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		// Parsear respuesta exitosa
		var apiResp usersAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if !apiResp.Success {
			return nil, fmt.Errorf("users-api returned success=false: %s", apiResp.Error)
		}

		// Decodificar el campo "data"
		var user User
		if err := json.Unmarshal(apiResp.Data, &user); err != nil {
			return nil, fmt.Errorf("failed to parse user data: %w", err)
		}

		return &user, nil

	case http.StatusNotFound:
		// Usuario no encontrado → retornar error de dominio
		return nil, domain.ErrDriverNotFound

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		// Error del servidor externo
		return nil, fmt.Errorf("users-api returned %d: service unavailable", resp.StatusCode)

	default:
		// Otros errores HTTP
		return nil, fmt.Errorf("users-api returned unexpected status %d: %s", resp.StatusCode, string(body))
	}
}
