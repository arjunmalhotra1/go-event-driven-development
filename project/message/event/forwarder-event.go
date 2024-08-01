package event

import "github.com/google/uuid"

type ForwarderEvent struct {
	BookingId       uuid.UUID `json:"bookingId"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
	ShowId          uuid.UUID `json:"showId"`
}
