package domain

type Event struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Venue  string `json:"venue"`
	Image  string `json:"image"`
	Status string `json:"status"`
}
