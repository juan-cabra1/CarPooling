package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TripDTO struct {
	ID             string    `json:"id"`
	DriverID       int64     `json:"driver_id"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	PricePerSeat   float64   `json:"price_per_seat"`
	Status         string    `json:"status"`
	DepartureDate  time.Time `json:"departure_date"`
}

type TripsClient interface {
	GetTrip(tripID string) (*TripDTO, error)
	UpdateTripSeats(tripID string, seatsChange int) error
	GetTripDriver(tripID string) (int64, error)
}

type tripsClientImpl struct {
	baseURL string
	client  *http.Client
}

func NewTripsClient(baseURL string) TripsClient {
	return &tripsClientImpl{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *tripsClientImpl) GetTrip(tripID string) (*TripDTO, error) {
	url := fmt.Sprintf("%s/trips/%s", c.baseURL, tripID)
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("trip not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Success bool     `json:"success"`
		Data    *TripDTO `json:"data"`
		Error   string   `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("api error: %s", response.Error)
	}

	return response.Data, nil
}

func (c *tripsClientImpl) UpdateTripSeats(tripID string, seatsChange int) error {
	url := fmt.Sprintf("%s/internal/trips/%s/seats", c.baseURL, tripID)
	
	payload := map[string]int{
		"seats_change": seatsChange,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update trip seats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update seats, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *tripsClientImpl) GetTripDriver(tripID string) (int64, error) {
	trip, err := c.GetTrip(tripID)
	if err != nil {
		return 0, err
	}
	return trip.DriverID, nil
}
