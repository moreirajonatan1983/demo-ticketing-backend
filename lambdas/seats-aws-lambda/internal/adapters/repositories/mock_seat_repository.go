package repositories

import (
	"fmt"

	"github.com/demoticketing/seats/internal/core/domain"
)

type MockSeatRepository struct{}

func NewMockSeatRepository() *MockSeatRepository {
	return &MockSeatRepository{}
}

func (r *MockSeatRepository) GetSeatsByEvent(eventId string) ([]domain.Seat, error) {
	// Generar una grilla mock (Filas A-H, Asientos 1-12)
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	var grid []domain.Seat

	for r := 0; r < 8; r++ {
		for s := 1; s <= 12; s++ {
			id := fmt.Sprintf("%s%d", rows[r], s)
			status := "available" // Default
			if r < 2 && s > 3 && s < 9 {
				status = "occupied"
			} else if r == 4 && s > 8 {
				status = "occupied"
			} else if r == 6 && (s == 2 || s == 3) {
				status = "processing"
			}

			grid = append(grid, domain.Seat{
				ID:     id,
				Row:    rows[r],
				Number: s,
				Status: status,
			})
		}
	}
	return grid, nil
}
