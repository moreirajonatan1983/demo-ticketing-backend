package services

import (
	"errors"
	"testing"

	"github.com/demoticketing/tickets/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTicketRepository
type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) GetTicketsByUser(userId string) ([]domain.Ticket, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketRepository) CreateTicket(ticket domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

// MockEventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) PublishTicketPurchased(ticket domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

func TestTicketService_GetTicketsForUser(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	mockPub := new(MockEventPublisher)
	service := NewTicketService(mockRepo, mockPub)

	expectedTickets := []domain.Ticket{
		{ID: "T1", UserID: "U1", Status: "confirmed"},
	}

	mockRepo.On("GetTicketsByUser", "U1").Return(expectedTickets, nil)

	tickets, err := service.GetTicketsForUser("U1")

	assert.NoError(t, err)
	assert.Equal(t, expectedTickets, tickets)
	mockRepo.AssertExpectations(t)
}

func TestTicketService_CreateTicket_Success(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	mockPub := new(MockEventPublisher)
	service := NewTicketService(mockRepo, mockPub)

	ticket := domain.Ticket{ID: "T1", UserID: "U1", Status: "confirmed"}

	mockRepo.On("CreateTicket", ticket).Return(nil)
	mockPub.On("PublishTicketPurchased", ticket).Return(nil)

	err := service.CreateTicket(ticket)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestTicketService_CreateTicket_RepoError(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	mockPub := new(MockEventPublisher)
	service := NewTicketService(mockRepo, mockPub)

	ticket := domain.Ticket{ID: "T1", UserID: "U1", Status: "confirmed"}

	mockRepo.On("CreateTicket", ticket).Return(errors.New("db error"))

	err := service.CreateTicket(ticket)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockRepo.AssertNotCalled(t, "PublishTicketPurchased", mock.Anything)
}

func TestTicketService_CreateTicket_PublisherErrorNonFatal(t *testing.T) {
	mockRepo := new(MockTicketRepository)
	mockPub := new(MockEventPublisher)
	service := NewTicketService(mockRepo, mockPub)

	ticket := domain.Ticket{ID: "T1", UserID: "U1", Status: "confirmed"}

	mockRepo.On("CreateTicket", ticket).Return(nil)
	mockPub.On("PublishTicketPurchased", ticket).Return(errors.New("sqs error"))

	err := service.CreateTicket(ticket)

	// Publisher error is non-fatal in the service (only logged)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}
