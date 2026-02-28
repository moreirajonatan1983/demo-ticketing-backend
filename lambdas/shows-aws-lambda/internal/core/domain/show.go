package domain

type Show struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Status  string `json:"status"` // available, soldout
}
