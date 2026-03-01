package services

import (
	"testing"

	"github.com/demoticketing/checkout/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestCheckoutService_ProcessCheckout_Approved(t *testing.T) {
	service := NewCheckoutService()
	req := domain.CheckoutRequest{
		EventId: "E1",
		Seats:   []string{"S1", "S2"},
		Method:  "CREDIT_CARD",
	}

	resp, err := service.ProcessCheckout(req)

	assert.NoError(t, err)
	assert.Equal(t, "APPROVED", resp.Status)
	assert.NotEmpty(t, resp.OrderID)
}

func TestCheckoutService_ProcessCheckout_Rejected(t *testing.T) {
	service := NewCheckoutService()
	req := domain.CheckoutRequest{
		EventId: "E1",
		Seats:   []string{"S1"},
		Method:  "REJECTED_CARD",
	}

	resp, err := service.ProcessCheckout(req)

	assert.NoError(t, err)
	assert.Equal(t, "REJECTED", resp.Status)
	assert.Empty(t, resp.OrderID)
}
