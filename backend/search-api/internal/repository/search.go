package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// SearchRepository handles search operations on trips
// This is a basic structure - full implementation will be in future phases
type SearchRepository interface {
	// Future methods will be added here for:
	// - SearchTrips (geospatial, text, filters)
	// - GetTripByID
	// - etc.
}

type searchRepository struct {
	collection *mongo.Collection
}

// NewSearchRepository creates a new search repository instance
func NewSearchRepository(db *mongo.Database) SearchRepository {
	return &searchRepository{
		collection: db.Collection("trips"),
	}
}
