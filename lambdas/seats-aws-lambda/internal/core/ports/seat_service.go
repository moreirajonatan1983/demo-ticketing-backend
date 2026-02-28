package ports

import "github.com/demoticketing/seats/internal/core/domain"

type SeatService interface {
	GetSeatsForEvent(eventId string) ([]domain.Seat, error)
}
