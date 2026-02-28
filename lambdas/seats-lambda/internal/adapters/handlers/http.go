package handlers

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/seats/internal/core/ports"
)

type HTTPHandler struct {
	service ports.SeatService
}

//go:embed docs/swagger.json
var docsJSON []byte

func NewHTTPHandler(service ports.SeatService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// @Summary Get seats for an event
// @Description Retrieve real-time seat status for a given event ID
// @Tags seats
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {array} domain.Seat
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events/{eventId}/seats [get]
func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	if request.Path == "/swagger.json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(docsJSON)}, nil
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
