package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) RemoveCanceledTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Removing cancelled tickets from the database")

	err := h.ticketsRepository.Delete(ctx, event.TicketID)
	if err != nil {
		return fmt.Errorf("failed to delete the ticket from the database %w", err)
	}
	return nil
}
