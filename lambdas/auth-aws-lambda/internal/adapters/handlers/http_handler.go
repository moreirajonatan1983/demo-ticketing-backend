package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/demoticketing/auth/internal/core/services"
)

// HTTPHandler despacha las rutas de la API de auth.
type HTTPHandler struct {
	svc *services.AuthService
}

func NewHTTPHandler(svc *services.AuthService) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

// HandleHTTPRequest es el entry point de la Lambda, mismo patrón que el resto.
func (h *HTTPHandler) HandleHTTPRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	corsHeaders := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "POST,GET,OPTIONS",
	}

	// Preflight CORS
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{StatusCode: 204, Headers: corsHeaders}, nil
	}

	switch {
	case req.HTTPMethod == "POST" && req.Path == "/auth/register":
		return h.handleRegister(req, corsHeaders)
	case req.HTTPMethod == "POST" && req.Path == "/auth/login":
		return h.handleLogin(req, corsHeaders)
	case req.HTTPMethod == "POST" && req.Path == "/auth/forgot-password":
		return h.handleForgotPassword(req, corsHeaders)
	case req.HTTPMethod == "GET" && req.Path == "/health":
		return jsonResponse(http.StatusOK, map[string]string{"status": "ok", "service": "auth"}, corsHeaders)
	default:
		return jsonResponse(http.StatusNotFound, map[string]string{"message": "Not found"}, corsHeaders)
	}
}

func (h *HTTPHandler) handleRegister(req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "Cuerpo inválido"}, headers)
	}

	result, err := h.svc.Register(services.RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			code = http.StatusConflict
		} else if errors.Is(err, services.ErrWeakPassword) {
			code = http.StatusBadRequest
		}
		return jsonResponse(code, map[string]string{"message": err.Error()}, headers)
	}

	return jsonResponse(http.StatusCreated, result, headers)
}

func (h *HTTPHandler) handleLogin(req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "Cuerpo inválido"}, headers)
	}

	result, err := h.svc.Login(services.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		code := http.StatusUnauthorized
		if errors.Is(err, services.ErrInvalidCredentials) {
			code = http.StatusUnauthorized
		}
		return jsonResponse(code, map[string]string{"message": err.Error()}, headers)
	}

	return jsonResponse(http.StatusOK, result, headers)
}

func (h *HTTPHandler) handleForgotPassword(req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return jsonResponse(http.StatusBadRequest, map[string]string{"message": "Cuerpo inválido"}, headers)
	}

	h.svc.ForgotPassword(input.Email)

	// Siempre responde 200 por seguridad
	return jsonResponse(http.StatusOK, map[string]string{
		"message": "Si el email existe, recibirás un enlace de recuperación.",
	}, headers)
}

func jsonResponse(statusCode int, body interface{}, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	b, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(b),
	}, nil
}
