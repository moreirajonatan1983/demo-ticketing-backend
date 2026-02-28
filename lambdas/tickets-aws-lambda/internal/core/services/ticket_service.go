package services

import (
	"github.com/demoticketing/tickets/internal/core/domain"
	"github.com/demoticketing/tickets/internal/core/ports"
)

type TicketService struct {
	repo ports.TicketRepository
}

func NewTicketService(repo ports.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) GetMyTickets(userId string) ([]domain.Ticket, error) {
	return s.repo.GetTicketsByUser(userId)
}
