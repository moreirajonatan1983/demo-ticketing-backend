package services

import (
	"log"

	"github.com/demoticketing/checkout/internal/core/domain"
)

type CheckoutService struct{}

func NewCheckoutService() *CheckoutService {
	return &CheckoutService{}
}

func (s *CheckoutService) ProcessCheckout(req domain.CheckoutRequest) (domain.CheckoutResponse, error) {
	log.Printf("Validating payment for Event: %s, Seats: %d, Method: %s", req.EventId, len(req.Seats), req.Method)

	if req.Method == "REJECTED_CARD" {
		return domain.CheckoutResponse{
			Message: "Pago Rechazado por Fondos Insuficientes",
			Status:  "REJECTED",
			OrderID: "",
		}, nil
	}

	return domain.CheckoutResponse{
		Message: "Pago exitoso",
		Status:  "APPROVED",
		OrderID: "ORD-XYZ-1234",
	}, nil
}
