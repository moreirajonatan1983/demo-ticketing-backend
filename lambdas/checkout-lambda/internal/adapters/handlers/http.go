package handlers

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/checkout/internal/core/domain"
	"github.com/demoticketing/checkout/internal/core/ports"
)

type HTTPHandler struct {
	service ports.CheckoutService
}

//go:embed docs/swagger.json
var docsJSON []byte

func NewHTTPHandler(service ports.CheckoutService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

// @Summary Process a checkout
// @Description Validates payload and triggers async queue for payment
// @Tags checkout
// @Accept json
// @Produce json
// @Param request body domain.CheckoutRequest true "Checkout Request object"
// @Success 202 {object} domain.CheckoutResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /checkout [post]
func (h *HTTPHandler) HandleHTTPRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Content-Type":                 "application/json",
	}

	if request.Path == "/swagger.json" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(docsJSON)}, nil
	}

	var req domain.CheckoutRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Headers: headers, Body: `{"error": "Invalid JSON"}`}, nil
	}

	res, err := h.service.ProcessCheckout(req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Processing Failed"}`}, nil
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
		Headers:    headers,
		Body:       string(body),
	}, nil
}
