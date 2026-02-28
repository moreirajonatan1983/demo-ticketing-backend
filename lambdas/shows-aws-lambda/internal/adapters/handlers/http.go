package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/shows/internal/core/ports"
)

type HTTPHandler struct {
	service ports.ShowService
}

func NewHTTPHandler(service ports.ShowService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) HandleHTTPRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventId := req.PathParameters["eventId"]

	if eventId == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing eventId",
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type",
			},
		}, nil
	}

	shows, err := h.service.GetShowsByEvent(eventId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error occurred",
		}, nil
	}

	bodyBytes, _ := json.Marshal(shows)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(bodyBytes),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}
