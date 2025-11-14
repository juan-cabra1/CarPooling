package repository

import (
	"context"
	"time"

	"search-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PopularRouteRepository handles popular route tracking
type PopularRouteRepository interface {
	IncrementSearchCount(ctx context.Context, originCity, destinationCity string) error
	GetTopRoutes(ctx context.Context, limit int) ([]domain.PopularRoute, error)
}

type popularRouteRepository struct {
	collection *mongo.Collection
}

// NewPopularRouteRepository creates a new popular route repository instance
func NewPopularRouteRepository(db *mongo.Database) PopularRouteRepository {
	return &popularRouteRepository{
		collection: db.Collection("popular_routes"),
	}
}

// IncrementSearchCount increments the search count for a route
// Uses upsert to create the route if it doesn't exist
func (r *popularRouteRepository) IncrementSearchCount(ctx context.Context, originCity, destinationCity string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{
		"origin_city":      originCity,
		"destination_city": destinationCity,
	}

	update := bson.M{
		"$inc": bson.M{
			"search_count": 1,
		},
		"$set": bson.M{
			"last_searched": time.Now(),
		},
		"$setOnInsert": bson.M{
			"origin_city":      originCity,
			"destination_city": destinationCity,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetTopRoutes returns the most popular routes
func (r *popularRouteRepository) GetTopRoutes(ctx context.Context, limit int) ([]domain.PopularRoute, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetSort(bson.D{{Key: "search_count", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var routes []domain.PopularRoute
	if err = cursor.All(ctx, &routes); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil
	if routes == nil {
		routes = []domain.PopularRoute{}
	}

	return routes, nil
}
