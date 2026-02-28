package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/seats/internal/core/ports"
)

type HTTPHandler struct {
	service ports.SeatService
}

func NewHTTPHandler(service ports.SeatService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	eventId, ok := request.PathParameters["eventId"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Headers: headers, Body: `{"error": "Missing eventId"}`}, nil
	}

	seats, err := h.service.GetSeatsForEvent(eventId)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Failed to get seats"}`}, nil
	}

	body, _ := json.Marshal(seats)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
