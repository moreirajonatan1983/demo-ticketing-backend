package ports

import "github.com/demoticketing/checkout/internal/core/domain"

type CheckoutService interface {
	ProcessCheckout(req domain.CheckoutRequest) (domain.CheckoutResponse, error)
}
