package ports

import "github.com/demoticketing/events/internal/core/domain"

type EventService interface {
	GetAllEvents() ([]domain.Event, error)
	GetEventByID(id string) (*domain.Event, error)
}
