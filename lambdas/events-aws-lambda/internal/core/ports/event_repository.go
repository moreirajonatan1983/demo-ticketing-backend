package ports

import "github.com/demoticketing/events/internal/core/domain"

type EventRepository interface {
	GetAll() ([]domain.Event, error)
	GetByID(id string) (*domain.Event, error)
}
