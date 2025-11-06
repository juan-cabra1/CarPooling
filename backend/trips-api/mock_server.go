package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Trip struct {
	ID             string    `json:"id"`
	DriverID       int64     `json:"driver_id"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	PricePerSeat   float64   `json:"price_per_seat"`
	Status         string    `json:"status"`
	DepartureDate  time.Time `json:"departure_date"`
}

var (
	trips = make(map[string]*Trip)
	mu    sync.RWMutex
)

func init() {
	// Inicializar con algunos viajes de prueba
	trips["trip001"] = &Trip{
		ID:             "trip001",
		DriverID:       3,
		TotalSeats:     4,
		AvailableSeats: 4,
		PricePerSeat:   50.00,
		Status:         "published",
		DepartureDate:  time.Now().Add(48 * time.Hour), // 2 días en el futuro
	}
	trips["trip002"] = &Trip{
		ID:             "trip002",
		DriverID:       4,
		TotalSeats:     3,
		AvailableSeats: 2,
		PricePerSeat:   75.00,
		Status:         "published",
		DepartureDate:  time.Now().Add(72 * time.Hour), // 3 días en el futuro
	}
	trips["trip003"] = &Trip{
		ID:             "trip003",
		DriverID:       3,
		TotalSeats:     5,
		AvailableSeats: 0,
		PricePerSeat:   60.00,
		Status:         "published",
		DepartureDate:  time.Now().Add(96 * time.Hour), // 4 días en el futuro
	}
	trips["trip004"] = &Trip{
		ID:             "trip004",
		DriverID:       4,
		TotalSeats:     4,
		AvailableSeats: 4,
		PricePerSeat:   100.00,
		Status:         "cancelled",
		DepartureDate:  time.Now().Add(120 * time.Hour), // 5 días en el futuro
	}
	trips["trip005"] = &Trip{
		ID:             "trip005",
		DriverID:       3,
		TotalSeats:     2,
		AvailableSeats: 2,
		PricePerSeat:   80.00,
		Status:         "published",
		DepartureDate:  time.Now().Add(24 * time.Hour), // 1 día en el futuro
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "trips-api mock is running",
		})
	})

	// Get trip by ID
	router.GET("/trips/:id", getTrip)

	// Update trip seats (internal endpoint)
	router.PUT("/internal/trips/:id/seats", updateTripSeats)

	// Create new trip (for testing)
	router.POST("/trips", createTrip)

	// List all trips (for debugging)
	router.GET("/trips", listTrips)

	log.Println("🚗 Trips API Mock server starting on port 8002...")
	if err := router.Run(":8002"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getTrip(c *gin.Context) {
	tripID := c.Param("id")

	mu.RLock()
	trip, exists := trips[tripID]
	mu.RUnlock()

	if !exists {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "trip not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    trip,
	})
}

func updateTripSeats(c *gin.Context) {
	tripID := c.Param("id")

	var req struct {
		SeatsChange int `json:"seats_change"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid request",
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	trip, exists := trips[tripID]
	if !exists {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "trip not found",
		})
		return
	}

	// Update available seats
	newAvailableSeats := trip.AvailableSeats + req.SeatsChange

	if newAvailableSeats < 0 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "not enough available seats",
		})
		return
	}

	if newAvailableSeats > trip.TotalSeats {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "cannot exceed total seats",
		})
		return
	}

	trip.AvailableSeats = newAvailableSeats

	log.Printf("Updated trip %s: seats_change=%d, new_available=%d", tripID, req.SeatsChange, newAvailableSeats)

	c.JSON(200, gin.H{
		"success": true,
		"data":    trip,
	})
}

func createTrip(c *gin.Context) {
	var req struct {
		DriverID      int64   `json:"driver_id" binding:"required"`
		TotalSeats    int     `json:"total_seats" binding:"required"`
		PricePerSeat  float64 `json:"price_per_seat" binding:"required"`
		DepartureDate string  `json:"departure_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid request: " + err.Error(),
		})
		return
	}

	departureDate, err := time.Parse(time.RFC3339, req.DepartureDate)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid departure_date format, use RFC3339",
		})
		return
	}

	mu.Lock()
	tripID := fmt.Sprintf("trip%03d", len(trips)+1)
	trip := &Trip{
		ID:             tripID,
		DriverID:       req.DriverID,
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.TotalSeats,
		PricePerSeat:   req.PricePerSeat,
		Status:         "active",
		DepartureDate:  departureDate,
	}
	trips[tripID] = trip
	mu.Unlock()

	log.Printf("Created new trip: %s (driver=%d, seats=%d)", tripID, req.DriverID, req.TotalSeats)

	c.JSON(201, gin.H{
		"success": true,
		"data":    trip,
	})
}

func listTrips(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	tripList := make([]*Trip, 0, len(trips))
	for _, trip := range trips {
		tripList = append(tripList, trip)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    tripList,
		"count":   len(tripList),
	})
}
