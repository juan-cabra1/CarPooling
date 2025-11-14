package repository

import (
	"context"
	"fmt"
	"time"
	"trips-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TripRepository define las operaciones de acceso a datos para viajes
type TripRepository interface {
	Create(ctx context.Context, trip *domain.Trip) error
	FindByID(ctx context.Context, id string) (*domain.Trip, error)
	FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error)
	Update(ctx context.Context, id string, trip *domain.Trip) error
	Delete(ctx context.Context, id string) error
	UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error
	Cancel(ctx context.Context, id string, cancelledBy int64, reason string) error
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Generar ObjectID si no existe
	if trip.ID.IsZero() {
		trip.ID = primitive.NewObjectID()
	}

	// Establecer timestamps
	now := time.Now()
	trip.CreatedAt = now
	trip.UpdatedAt = now

	// Inicializar availability_version si es 0
	if trip.AvailabilityVersion == 0 {
		trip.AvailabilityVersion = 1
	}

	_, err := r.collection.InsertOne(ctx, trip)
	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

// FindByID busca un viaje por su ID
func (r *tripRepository) FindByID(ctx context.Context, id string) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid trip ID format: %w", err)
	}

	var trip domain.Trip
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&trip)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTripNotFound
		}
		return nil, fmt.Errorf("failed to find trip: %w", err)
	}

	return &trip, nil
}

// FindAll busca viajes con filtros y paginación
func (r *tripRepository) FindAll(ctx context.Context, filters map[string]interface{}, page, limit int) ([]domain.Trip, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Construir filtro MongoDB
	filter := bson.M{}
	for key, value := range filters {
		filter[key] = value
	}

	// Contar total de documentos
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count trips: %w", err)
	}

	// Calcular skip para paginación
	skip := (page - 1) * limit

	// Opciones de búsqueda con paginación
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "departure_datetime", Value: 1}}) // Ordenar por fecha de salida ascendente

	// Ejecutar búsqueda
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trips: %w", err)
	}
	defer cursor.Close(ctx)

	// Decodificar resultados
	var trips []domain.Trip
	if err = cursor.All(ctx, &trips); err != nil {
		return nil, 0, fmt.Errorf("failed to decode trips: %w", err)
	}

	// Si no hay resultados, retornar slice vacío en lugar de nil
	if trips == nil {
		trips = []domain.Trip{}
	}

	return trips, total, nil
}

// Update actualiza un viaje existente
func (r *tripRepository) Update(ctx context.Context, id string, trip *domain.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid trip ID format: %w", err)
	}

	// Actualizar updated_at
	trip.UpdatedAt = time.Now()

	// Usar $set para actualizar solo los campos proporcionados
	update := bson.M{
		"$set": trip,
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

// Delete elimina un viaje (hard delete)
func (r *tripRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid trip ID format: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

// UpdateAvailability actualiza la disponibilidad de asientos con optimistic locking
// CRÍTICO: Este método previene race conditions usando availability_version
func (r *tripRepository) UpdateAvailability(ctx context.Context, tripID string, seatsDelta int, expectedVersion int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return fmt.Errorf("invalid trip ID format: %w", err)
	}

	// Filtro con optimistic locking
	// Solo actualiza si:
	// 1. El trip existe
	// 2. La versión coincide (expectedVersion)
	// 3. Hay suficientes asientos disponibles
	filter := bson.M{
		"_id":                  objectID,
		"availability_version": expectedVersion,
		"available_seats":      bson.M{"$gte": -seatsDelta}, // Si seatsDelta es negativo (reserva), available_seats debe ser >= abs(seatsDelta)
	}

	// Update atómico que incrementa/decrementa asientos y actualiza la versión
	update := bson.M{
		"$inc": bson.M{
			"available_seats":      seatsDelta,      // +N para cancelaciones, -N para reservas
			"reserved_seats":       -seatsDelta,     // -N para cancelaciones, +N para reservas
			"availability_version": 1,               // Siempre incrementar la versión
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update availability: %w", err)
	}

	// Si MatchedCount es 0, significa que:
	// - El trip no existe, O
	// - La versión no coincide (conflicto de concurrencia), O
	// - No hay suficientes asientos disponibles
	if result.MatchedCount == 0 {
		return domain.ErrOptimisticLockFailed
	}

	return nil
}

// Cancel marca un viaje como cancelado
func (r *tripRepository) Cancel(ctx context.Context, id string, cancelledBy int64, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid trip ID format: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":              "cancelled",
			"cancelled_at":        &now,
			"cancelled_by":        &cancelledBy,
			"cancellation_reason": reason,
			"updated_at":          now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to cancel trip: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}
