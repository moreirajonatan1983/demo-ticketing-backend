package repositories

import "github.com/demoticketing/shows/internal/core/domain"

type mockShowRepository struct{}

func NewMockShowRepository() *mockShowRepository {
	return &mockShowRepository{}
}

func (m *mockShowRepository) GetShowsByEvent(eventId string) ([]domain.Show, error) {
	return []domain.Show{
		{ID: "1", EventID: eventId, Date: "15 Oct 2026", Time: "21:00", Status: "available"},
		{ID: "2", EventID: eventId, Date: "16 Oct 2026", Time: "20:00", Status: "soldout"},
	}, nil
}
