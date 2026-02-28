package domain

type Ticket struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	EventName string `json:"event_name"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	Sector    string `json:"sector"`
	Status    string `json:"status"` // PAID, PROCESSING
}
