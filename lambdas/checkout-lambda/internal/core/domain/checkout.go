package domain

type CheckoutRequest struct {
	EventId string   `json:"event_id"`
	Seats   []string `json:"seats"`
	Method  string   `json:"method"`
}

type CheckoutResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}
