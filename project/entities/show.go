package entities

type Show struct {
	ShowId          string `json:"show_id" db:"show_id"`
	DeadNationId    string `json:"dead_nation_id" db:"dead_nation_id"`
	NumberOfTickets int64  `json:"number_of_tickets" db:"number_of_tickets"`
	StartTime       string `json:"start_time" db:"start_time"`
	Title           string `json:"title" db:"title"`
	Venue           string `json:"venue" db:"venue"`
}
