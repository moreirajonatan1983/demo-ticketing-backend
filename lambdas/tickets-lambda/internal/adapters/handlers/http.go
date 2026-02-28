package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/tickets/internal/core/ports"
)

type HTTPHandler struct {
	service ports.TicketService
}

func NewHTTPHandler(service ports.TicketService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	// Assuming a mock logged in user ID via arbitrary query parameter for now
	userId := request.QueryStringParameters["userId"]
	if userId == "" {
		userId = "mock-user-123" // Fallback to mock user if nothing is passed
	}

	tickets, err := h.service.GetMyTickets(userId)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Failed to get tickets"}`}, nil
	}

	body, _ := json.Marshal(tickets)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
