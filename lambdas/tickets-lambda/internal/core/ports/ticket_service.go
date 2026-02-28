package ports

import "github.com/demoticketing/tickets/internal/core/domain"

type TicketService interface {
	GetMyTickets(userId string) ([]domain.Ticket, error)
}
