package services

import (
	"github.com/demoticketing/events/internal/core/domain"
	"github.com/demoticketing/events/internal/core/ports"
)

type EventService struct {
	repo ports.EventRepository
}

func NewEventService(repo ports.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAllEvents() ([]domain.Event, error) {
	return s.repo.GetAll()
}

func (s *EventService) GetEventByID(id string) (*domain.Event, error) {
	return s.repo.GetByID(id)
}
