package services

import (
	"github.com/demoticketing/seats/internal/core/domain"
	"github.com/demoticketing/seats/internal/core/ports"
)

type SeatService struct {
	repo ports.SeatRepository
}

func NewSeatService(repo ports.SeatRepository) *SeatService {
	return &SeatService{repo: repo}
}

func (s *SeatService) GetSeatsForEvent(eventId string) ([]domain.Seat, error) {
	return s.repo.GetSeatsByEvent(eventId)
}
