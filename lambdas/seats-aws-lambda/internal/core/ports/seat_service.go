package ports

import "github.com/demoticketing/seats/internal/core/domain"

type SeatService interface {
	GetSeatsForEvent(eventId string) ([]domain.Seat, error)
	ReserveSeat(eventId, seatId string) error
	ReleaseSeat(eventId, seatId string) error
}
