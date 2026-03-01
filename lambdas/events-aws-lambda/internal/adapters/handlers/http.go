package handlers

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/events/internal/core/ports"
)

type HTTPHandler struct {
	service ports.EventService
}

//go:embed docs/swagger.json
var docsJSON []byte

func NewHTTPHandler(service ports.EventService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// @Summary Get event(s)
// @Description Retrieve events list or specific event by ID
// @Tags events
// @Produce json
// @Param id path string false "Event ID"
// @Success 200 {object} interface{}
// @Failure 404 {object} map[string]string
// @Router /events [get]
// @Router /events/{id} [get]
func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	if request.Path == "/swagger.json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(docsJSON)}, nil
	}

	id, hasId := request.PathParameters["id"]

	// B-20: GET /events/{id}/limits — Devuelve restricciones de compra por usuario
	if hasId && strings.HasSuffix(request.Path, "/limits") {
		limits, err := h.service.GetEventLimits(id)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Headers: headers, Body: `{"error": "Event not found"}`}, nil
		}
		body, _ := json.Marshal(limits)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil
	}

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
