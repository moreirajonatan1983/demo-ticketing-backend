package ports

import "github.com/demoticketing/shows/internal/core/domain"

type ShowService interface {
	GetShowsByEvent(eventId string) ([]domain.Show, error)
}
