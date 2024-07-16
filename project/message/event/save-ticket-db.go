package event

import (
	"context"
	"fmt"
	"tickets/entities"
)

func (h Handler) SaveTicketToDatabase(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	err := h.ticketsRepository.Add(ctx, entities.Ticket{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	})

	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}
