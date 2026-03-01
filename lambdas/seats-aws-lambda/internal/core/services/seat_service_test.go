package services

import (
	"errors"
	"testing"

	"github.com/demoticketing/seats/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeatRepository is a mock implementation of the SeatRepository interface
type MockSeatRepository struct {
	mock.Mock
}

func (m *MockSeatRepository) GetSeatsByEvent(eventId string) ([]domain.Seat, error) {
	args := m.Called(eventId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Seat), args.Error(1)
}

func (m *MockSeatRepository) ReserveSeat(eventId string, seatId string) error {
	args := m.Called(eventId, seatId)
	return args.Error(0)
}

func (m *MockSeatRepository) ReleaseSeat(eventId string, seatId string) error {
	args := m.Called(eventId, seatId)
	return args.Error(0)
}

func TestSeatService_GetSeatsForEvent(t *testing.T) {
	mockRepo := new(MockSeatRepository)
	service := NewSeatService(mockRepo)

	expectedSeats := []domain.Seat{
		{ID: "S1", Row: "A", Number: 1, Status: "available"},
		{ID: "S2", Row: "A", Number: 2, Status: "reserved"},
	}

	mockRepo.On("GetSeatsByEvent", "E1").Return(expectedSeats, nil)

	seats, err := service.GetSeatsForEvent("E1")

	assert.NoError(t, err)
	assert.Equal(t, expectedSeats, seats)
	mockRepo.AssertExpectations(t)
}

func TestSeatService_ReserveSeat_Success(t *testing.T) {
	mockRepo := new(MockSeatRepository)
	service := NewSeatService(mockRepo)

	mockRepo.On("ReserveSeat", "E1", "S1").Return(nil)

	err := service.ReserveSeat("E1", "S1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSeatService_ReserveSeat_Error(t *testing.T) {
	mockRepo := new(MockSeatRepository)
	service := NewSeatService(mockRepo)

	mockRepo.On("ReserveSeat", "E1", "S1").Return(errors.New("db error"))

	err := service.ReserveSeat("E1", "S1")

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestSeatService_ReleaseSeat(t *testing.T) {
	mockRepo := new(MockSeatRepository)
	service := NewSeatService(mockRepo)

	mockRepo.On("ReleaseSeat", "E1", "S1").Return(nil)

	err := service.ReleaseSeat("E1", "S1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
