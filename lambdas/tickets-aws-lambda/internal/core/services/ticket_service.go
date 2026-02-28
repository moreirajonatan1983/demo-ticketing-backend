package services

import (
	"log"

	"github.com/demoticketing/tickets/internal/core/domain"
	"github.com/demoticketing/tickets/internal/core/ports"
)

type TicketService struct {
	repo      ports.TicketRepository
	publisher ports.EventPublisher
}

func NewTicketService(repo ports.TicketRepository, publisher ports.EventPublisher) *TicketService {
	return &TicketService{repo: repo, publisher: publisher}
}

func (s *TicketService) GetTicketsForUser(userId string) ([]domain.Ticket, error) {
	return s.repo.GetTicketsByUser(userId)
}

func (s *TicketService) CreateTicket(ticket domain.Ticket) error {
	if err := s.repo.CreateTicket(ticket); err != nil {
		return err
	}
	// B-12: Emit TicketPurchased event to SQS so ticket-worker and notification-service react
	if err := s.publisher.PublishTicketPurchased(ticket); err != nil {
		log.Printf("[WARNING] Ticket created but TicketPurchased event failed to publish: %v", err)
		// Non-fatal: ticket is already in DynamoDB
	}
	return nil
}
