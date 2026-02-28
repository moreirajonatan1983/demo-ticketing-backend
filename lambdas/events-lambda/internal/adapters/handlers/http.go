package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/events/internal/core/ports"
)

type HTTPHandler struct {
	service ports.EventService
}

func NewHTTPHandler(service ports.EventService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	id, hasId := request.PathParameters["id"]

	var body []byte
	var err error

	if hasId {
		event, errServ := h.service.GetEventByID(id)
		if errServ != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Headers: headers, Body: `{"error": "Event not found"}`}, nil
		}
		body, err = json.Marshal(event)
	} else {
		allEvents, _ := h.service.GetAllEvents()
		body, err = json.Marshal(allEvents)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Failed to marshal response"}`}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
