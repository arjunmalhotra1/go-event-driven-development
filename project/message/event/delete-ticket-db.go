package event

import (
	"context"
	"fmt"
	"tickets/entities"
)

func (h Handler) Delete(ctx context.Context, event *entities.TicketBookingCanceled) error {
	err := h.ticketsRepository.Delete(ctx, event.TicketID)
	if err != nil {
		return fmt.Errorf("failed to delete the ticket from the database %w", err)
	}
	return nil
}
