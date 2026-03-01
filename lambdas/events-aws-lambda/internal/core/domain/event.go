package domain

type Event struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Date              string `json:"date"`
	Venue             string `json:"venue"`
	Image             string `json:"image"`
	Status            string `json:"status"`
	Description       string `json:"description,omitempty"`
	MaxTicketsPerUser int    `json:"max_tickets_per_user,omitempty"`
}

// EventLimits representa las restricciones de compra de un evento
type EventLimits struct {
	EventID           string `json:"event_id"`
	MaxTicketsPerUser int    `json:"max_tickets_per_user"`
	Message           string `json:"message"`
}
