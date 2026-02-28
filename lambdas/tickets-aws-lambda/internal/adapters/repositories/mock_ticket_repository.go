package repositories

import "github.com/demoticketing/tickets/internal/core/domain"

type MockTicketRepository struct{}

func NewMockTicketRepository() *MockTicketRepository {
	return &MockTicketRepository{}
}

func (r *MockTicketRepository) GetTicketsByUser(userId string) ([]domain.Ticket, error) {
	// Return a mock user ticket list
	return []domain.Ticket{
		{
			ID:        "TCK-9902-ABCD",
			EventName: "COLDPLAY - Music of the Spheres",
			Date:      "15 Octubre 2026",
			Location:  "Estadio Nacional",
			Sector:    "VIP Front Stage",
			Status:    "PAGADO OKEY",
		},
	}, nil
}
