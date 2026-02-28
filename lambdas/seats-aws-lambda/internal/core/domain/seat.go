package domain

type Seat struct {
	ID     string `json:"id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
	Status string `json:"status"` // available, occupied, processing
}
