package controller

import (
	"net/http"
	"strconv"
	"strings"
	"trips-api/internal/domain"
	"trips-api/internal/service"

	"github.com/gin-gonic/gin"
)

// TripController define la interfaz del controlador de viajes
type TripController interface {
	CreateTrip(c *gin.Context)
	GetTrip(c *gin.Context)
	ListTrips(c *gin.Context)
	UpdateTrip(c *gin.Context)
	DeleteTrip(c *gin.Context)
}

type tripController struct {
	tripService service.TripService
}

// NewTripController crea una nueva instancia del controlador de viajes
func NewTripController(tripService service.TripService) TripController {
	return &tripController{
		tripService: tripService,
	}
}

// CreateTrip maneja la creación de un nuevo viaje
// POST /trips
// Requiere autenticación (JWT)
func (ctrl *tripController) CreateTrip(c *gin.Context) {
	// Extraer user_id del contexto (viene del middleware JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Extraer Authorization header para forwarding a users-api
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "token de autenticación requerido",
		})
		return
	}

	// Validar formato del header (debe ser "Bearer {token}")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "formato de token inválido",
		})
		return
	}

	// Bind request body a CreateTripRequest
	var request domain.CreateTripRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	// Llamar al servicio (forward complete Authorization header)
	trip, err := ctrl.tripService.CreateTrip(c.Request.Context(), userID.(int64), authHeader, request)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(201, gin.H{
		"success": true,
		"data":    trip,
	})
}

// GetTrip obtiene un viaje por su ID
// GET /trips/:id
// Público (sin autenticación)
func (ctrl *tripController) GetTrip(c *gin.Context) {
	// Extraer trip ID del path
	tripID := c.Param("id")

	// Llamar al servicio
	trip, err := ctrl.tripService.GetTrip(c.Request.Context(), tripID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(200, gin.H{
		"success": true,
		"data":    trip,
	})
}

// ListTrips lista viajes con filtros y paginación
// GET /trips?driver_id=X&status=Y&origin_city=Z&destination_city=W&page=1&limit=10
// Público (sin autenticación)
func (ctrl *tripController) ListTrips(c *gin.Context) {
	// Extraer query parameters
	driverIDStr := c.Query("driver_id")
	status := c.Query("status")
	originCity := c.Query("origin_city")
	destinationCity := c.Query("destination_city")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parsear paginación
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Máximo 100 items por página
	}

	// Construir map de filtros
	filters := make(map[string]interface{})

	if driverIDStr != "" {
		driverID, err := strconv.ParseInt(driverIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "driver_id inválido",
			})
			return
		}
		filters["driver_id"] = driverID
	}

	if status != "" {
		filters["status"] = status
	}

	if originCity != "" {
		filters["origin.city"] = originCity
	}

	if destinationCity != "" {
		filters["destination.city"] = destinationCity
	}

	// Llamar al servicio
	trips, total, err := ctrl.tripService.ListTrips(c.Request.Context(), filters, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Respuesta exitosa con paginación
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"trips": trips,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// UpdateTrip actualiza un viaje existente
// PUT /trips/:id
// Requiere autenticación (JWT) y ser el dueño del viaje
func (ctrl *tripController) UpdateTrip(c *gin.Context) {
	// Extraer trip ID del path
	tripID := c.Param("id")

	// Extraer user_id del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Extraer role del contexto (viene del middleware JWT)
	userRole, roleExists := c.Get("role")
	if !roleExists {
		userRole = "user" // default
	}

	// Bind request body a UpdateTripRequest
	var request domain.UpdateTripRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	// Llamar al servicio (el servicio valida ownership)
	trip, err := ctrl.tripService.UpdateTrip(c.Request.Context(), tripID, userID.(int64), userRole.(string), request)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(200, gin.H{
		"success": true,
		"data":    trip,
	})
}

// DeleteTrip elimina un viaje
// DELETE /trips/:id
// Requiere autenticación (JWT) y ser el dueño del viaje
func (ctrl *tripController) DeleteTrip(c *gin.Context) {
	// Extraer trip ID del path
	tripID := c.Param("id")

	// Extraer user_id del contexto
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	// Extraer role del contexto (viene del middleware JWT)
	userRole, roleExists := c.Get("role")
	if !roleExists {
		userRole = "user" // default
	}

	// Llamar al servicio (el servicio valida ownership)
	err := ctrl.tripService.DeleteTrip(c.Request.Context(), tripID, userID.(int64), userRole.(string))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Respuesta exitosa
	c.JSON(200, gin.H{
		"success": true,
		"message": "trip deleted successfully",
	})
}

// handleServiceError maneja los errores del servicio y los mapea a status codes HTTP
func handleServiceError(c *gin.Context, err error) {
	// Type assertion a AppError
	if appErr, ok := err.(*domain.AppError); ok {
		switch appErr.Code {
		case "TRIP_NOT_FOUND", "DRIVER_NOT_FOUND":
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   appErr.Message,
			})
		case "UNAUTHORIZED":
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   appErr.Message,
			})
		case "PAST_DEPARTURE", "HAS_RESERVATIONS", "NO_SEATS_AVAILABLE":
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   appErr.Message,
			})
		case "OPTIMISTIC_LOCK_FAILED":
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   appErr.Message,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   appErr.Message,
			})
		}
		return
	}

	// Error genérico
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
