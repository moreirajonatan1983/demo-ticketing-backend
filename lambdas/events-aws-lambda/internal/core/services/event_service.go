package services

import (
	"fmt"

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

// GetEventLimits retorna las restricciones de compra de un evento (B-20)
func (s *EventService) GetEventLimits(id string) (*domain.EventLimits, error) {
	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch event: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("event not found: %s", id)
	}

	maxPerUser := event.MaxTicketsPerUser
	if maxPerUser == 0 {
		maxPerUser = 4 // Límite por defecto
	}

	return &domain.EventLimits{
		EventID:           event.ID,
		MaxTicketsPerUser: maxPerUser,
		Message:           fmt.Sprintf("Podés comprar hasta %d entradas por usuario para este evento.", maxPerUser),
	}, nil
}
