package ports

import "github.com/demoticketing/shows/internal/core/domain"

type ShowRepository interface {
	GetShowsByEvent(eventId string) ([]domain.Show, error)
}
