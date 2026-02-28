package ports

import "github.com/demoticketing/tickets/internal/core/domain"

type TicketRepository interface {
	GetTicketsByUser(userId string) ([]domain.Ticket, error)
}
