package handlers

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/tickets/internal/core/domain"
	"github.com/demoticketing/tickets/internal/core/ports"
)

type HTTPHandler struct {
	service ports.TicketService
	storage ports.TicketStorage
}

//go:embed docs/swagger.json
var docsJSON []byte

func NewHTTPHandler(service ports.TicketService, storage ports.TicketStorage) *HTTPHandler {
	return &HTTPHandler{service: service, storage: storage}
}

// @Summary Get user tickets
// @Description Retrieve real-time purchased tickets for the logged user
// @Tags tickets
// @Produce json
// @Param userId query string false "User ID"
// @Success 200 {array} domain.Ticket
// @Failure 500 {object} map[string]string
// @Router /tickets/me [get]
func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	if request.Path == "/swagger.json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(docsJSON)}, nil
	}

	// POST /tickets handler
	if request.HTTPMethod == "POST" && request.Resource == "/tickets" {
		var ticket domain.Ticket
		if err := json.Unmarshal([]byte(request.Body), &ticket); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Headers: headers, Body: `{"error": "Invalid JSON"}`}, nil
		}

		err := h.service.CreateTicket(ticket)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Failed to create ticket"}`}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: http.StatusCreated, Headers: headers, Body: `{"message": "Ticket created"}`}, nil
	}

	// GET /tickets/me handler
	if request.HTTPMethod == "GET" && request.Resource == "/tickets/me" {
		userId := request.QueryStringParameters["userId"]
		if userId == "" {
			userId = "mock-user-123" // Fallback to mock user if nothing is passed
		}

		tickets, err := h.service.GetTicketsForUser(userId)
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

	// GET /tickets/{ticketId}/download → presigned S3 URL for PDF
	if request.HTTPMethod == "GET" && len(request.PathParameters) > 0 {
		ticketId, ok := request.PathParameters["ticketId"]
		if ok && request.Resource == "/tickets/{ticketId}/download" {
			url, err := h.storage.GetPresignedDownloadURL(ticketId)
			if err != nil {
				body, _ := json.Marshal(map[string]string{"error": "PDF not found or not ready yet"})
				return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Headers: headers, Body: string(body)}, nil
			}
			body, _ := json.Marshal(map[string]string{"downloadUrl": url})
			return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil
		}
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed, Headers: headers, Body: `{"error": "Method not allowed"}`}, nil
}
