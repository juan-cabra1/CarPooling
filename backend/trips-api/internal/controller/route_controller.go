package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouteController maneja el endpoint de rutas (Google Maps Embed)
type RouteController interface {
	GetRouteEmbed(c *gin.Context)
}

type routeController struct {
	googleMapsAPIKey string
}

func NewRouteController(googleMapsAPIKey string) RouteController {
	return &routeController{googleMapsAPIKey: googleMapsAPIKey}
}

// GetRouteEmbed devuelve la URL del iframe de Google Maps Directions.
// El API key se mantiene en el servidor y nunca se expone en el bundle del frontend.
//
// GET /route?origin_lat=...&origin_lng=...&dest_lat=...&dest_lng=...
// Respuesta: { "embed_url": "https://www.google.com/maps/embed/v1/directions?..." }
func (rc *routeController) GetRouteEmbed(c *gin.Context) {
	if rc.googleMapsAPIKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "maps no configurado",
		})
		return
	}

	originLat := c.Query("origin_lat")
	originLng := c.Query("origin_lng")
	destLat := c.Query("dest_lat")
	destLng := c.Query("dest_lng")

	// Fallback a nombres de ciudad si no hay coordenadas
	originCity := c.Query("origin_city")
	destCity := c.Query("dest_city")
	originProvince := c.Query("origin_province")
	destProvince := c.Query("dest_province")

	var origin, destination string

	if originLat != "" && originLng != "" {
		origin = fmt.Sprintf("%s,%s", originLat, originLng)
	} else if originCity != "" {
		if originProvince != "" {
			origin = fmt.Sprintf("%s, %s, Argentina", originCity, originProvince)
		} else {
			origin = fmt.Sprintf("%s, Argentina", originCity)
		}
	}

	if destLat != "" && destLng != "" {
		destination = fmt.Sprintf("%s,%s", destLat, destLng)
	} else if destCity != "" {
		if destProvince != "" {
			destination = fmt.Sprintf("%s, %s, Argentina", destCity, destProvince)
		} else {
			destination = fmt.Sprintf("%s, Argentina", destCity)
		}
	}

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "se requieren origen y destino (lat/lng o ciudad)",
		})
		return
	}

	embedURL := fmt.Sprintf(
		"https://www.google.com/maps/embed/v1/directions?key=%s&origin=%s&destination=%s&mode=driving&language=es",
		rc.googleMapsAPIKey,
		origin,
		destination,
	)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"embed_url": embedURL,
	})
}
