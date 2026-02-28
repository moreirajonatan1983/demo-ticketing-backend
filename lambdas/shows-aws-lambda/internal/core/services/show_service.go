package services

import (
	"github.com/demoticketing/shows/internal/core/domain"
	"github.com/demoticketing/shows/internal/core/ports"
)

type showService struct {
	repo ports.ShowRepository
}

func NewShowService(repo ports.ShowRepository) ports.ShowService {
	return &showService{repo: repo}
}

func (s *showService) GetShowsByEvent(eventId string) ([]domain.Show, error) {
	return s.repo.GetShowsByEvent(eventId)
}
