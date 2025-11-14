package mocks

import (
	"github.com/juan-cabra1/CarPooling/backend/search-api/internal/domain"
)

// MockSolrClient is a mock implementation of the Solr client
type MockSolrClient struct {
	IndexFunc   func(trip *domain.SearchTrip) error
	DeleteFunc  func(tripID string) error
	SearchFunc  func(query map[string]interface{}) ([]domain.SearchTrip, int64, error)
	PingFunc    func() error
	IsAvailable bool
}

// Index calls the mocked IndexFunc
func (m *MockSolrClient) Index(trip *domain.SearchTrip) error {
	if m.IndexFunc != nil {
		return m.IndexFunc(trip)
	}
	return nil
}

// Delete calls the mocked DeleteFunc
func (m *MockSolrClient) Delete(tripID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(tripID)
	}
	return nil
}

// Search calls the mocked SearchFunc
func (m *MockSolrClient) Search(query map[string]interface{}) ([]domain.SearchTrip, int64, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(query)
	}
	return []domain.SearchTrip{}, 0, nil
}

// Ping calls the mocked PingFunc
func (m *MockSolrClient) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}
