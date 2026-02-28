package ports

import "github.com/demoticketing/tickets/internal/core/domain"

type TicketService interface {
	GetTicketsForUser(userId string) ([]domain.Ticket, error)
	CreateTicket(ticket domain.Ticket) error
}
