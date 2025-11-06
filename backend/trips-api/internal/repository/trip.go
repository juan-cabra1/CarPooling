package repository

import (
	"context"
	"trips-api/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

// TripRepository define las operaciones de acceso a datos para viajes
type TripRepository interface {
	Create(ctx context.Context, trip *domain.Trip) error
	FindByID(ctx context.Context, id string) (*domain.Trip, error)
	FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error)
	Update(ctx context.Context, id string, trip *domain.Trip) error
	Delete(ctx context.Context, id string) error
	UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error
}

type tripRepository struct {
	collection *mongo.Collection
}

// NewTripRepository crea una nueva instancia del repositorio de viajes
func NewTripRepository(db *mongo.Database) TripRepository {
	return &tripRepository{
		collection: db.Collection("trips"),
	}
}

// Create inserta un nuevo viaje en la base de datos
func (r *tripRepository) Create(ctx context.Context, trip *domain.Trip) error {
	// TODO: Implementar en Phase 3
	return nil
}

// FindByID busca un viaje por su ID
func (r *tripRepository) FindByID(ctx context.Context, id string) (*domain.Trip, error) {
	// TODO: Implementar en Phase 3
	return nil, nil
}

// FindAll busca viajes con filtros y paginación
func (r *tripRepository) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error) {
	// TODO: Implementar en Phase 3
	return nil, 0, nil
}

// Update actualiza un viaje existente
func (r *tripRepository) Update(ctx context.Context, id string, trip *domain.Trip) error {
	// TODO: Implementar en Phase 3
	return nil
}

// Delete elimina un viaje (soft delete)
func (r *tripRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implementar en Phase 3
	return nil
}

// UpdateAvailability actualiza la disponibilidad de asientos con optimistic locking
func (r *tripRepository) UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error {
	// TODO: Implementar en Phase 3 - CRÍTICO para evitar race conditions
	// Debe usar availability_version para optimistic locking
	return nil
}
