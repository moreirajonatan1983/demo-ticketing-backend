package ports

import "github.com/demoticketing/seats/internal/core/domain"

type SeatRepository interface {
	GetSeatsByEvent(eventId string) ([]domain.Seat, error)
	ReserveSeat(eventId string, seatId string) error
	ReleaseSeat(eventId string, seatId string) error
}
