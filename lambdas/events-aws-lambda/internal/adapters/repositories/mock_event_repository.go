package repositories

import (
	"errors"

	"github.com/demoticketing/events/internal/core/domain"
)

type MockEventRepository struct {
	events []domain.Event
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events: []domain.Event{
			{ID: "1", Title: "COLDPLAY - Music of the Spheres", Date: "15 OCT 2026", Venue: "Estadio Nacional", Image: "https://images.unsplash.com/photo-1540039155733-d7696d819924?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80", Status: "Disponible"},
			{ID: "2", Title: "THE WEEKND - After Hours Tour", Date: "22 NOV 2026", Venue: "Movistar Arena", Image: "https://images.unsplash.com/photo-1493225457124-a1a2a5f5f4b2?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80", Status: "Pocos Tickets"},
			{ID: "3", Title: "DUA LIPA - Radical Optimism", Date: "04 DIC 2026", Venue: "Hipódromo", Image: "https://images.unsplash.com/photo-1514525253161-7a46d19cd819?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80", Status: "Sold Out"},
		},
	}
}

func (r *MockEventRepository) GetAll() ([]domain.Event, error) {
	return r.events, nil
}

func (r *MockEventRepository) GetByID(id string) (*domain.Event, error) {
	for _, e := range r.events {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, errors.New("event not found")
}
